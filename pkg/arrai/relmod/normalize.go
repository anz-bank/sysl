//nolint:dupl
package relmod

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

// tuple is a type alias for the ugly common map type.
type tuple map[string]interface{}

// Normalize transforms a module into a relational model schema.
func Normalize(m *sysl.Module) (*Schema, error) {
	s := &Schema{}
	if err := normalizeModule(s, m); err != nil {
		return nil, err
	}
	return s, nil
}

func normalizeModule(s *Schema, m *sysl.Module) error {
	for _, app := range m.Apps {
		if err := normalizeApp(s, app); err != nil {
			return err
		}
	}
	return nil
}

func normalizeApp(s *Schema, app *sysl.Application) error {
	s.App = append(s.App, App{
		AppName:      app.Name.Part,
		AppLongName:  app.LongName,
		AppDocstring: app.Docstring,
	})

	normalizeAppMeta(s, app)

	for _, mixin := range app.Mixin2 {
		normalizeMixin(s, app, mixin)
	}

	for _, ep := range app.Endpoints {
		if err := normalizeEndpoint(s, app, ep); err != nil {
			return err
		}
	}

	for typeName, typ := range app.Types {
		normalizeType(s, app, typ, typeName)
	}

	for viewName, view := range app.Views {
		normalizeView(s, app, view, viewName)
	}

	return nil
}

func normalizeMixin(s *Schema, app *sysl.Application, mixin *sysl.Application) {
	s.Mixin = append(s.Mixin, Mixin{
		AppName:   app.Name.Part,
		MixinName: mixin.Name.Part,
	})

	normalizeMixinMeta(s, app, mixin)
}

func normalizeEndpoint(s *Schema, app *sysl.Application, ep *sysl.Endpoint) error {
	if ep.Name == placeholder {
		return nil
	}

	if ep.IsPubsub {
		normalizeEvent(s, app, ep)
		return nil
	}

	var epEvent EndpointEvent
	if ep.Source != nil {
		epEvent = EndpointEvent{
			AppName:   EndpointEventAppName{Part: ep.Source.Part},
			EventName: strings.SplitAfter(ep.Name, " -> ")[1],
		}
	}

	rest, err := parseRestPath(ep.RestParams)
	if err != nil {
		return err
	}

	s.Ep = append(s.Ep, Endpoint{
		AppName:     app.Name.Part,
		EpName:      ep.Name,
		EpLongName:  ep.LongName,
		EpDocstring: ep.Docstring,
		EpEvent:     epEvent,
		Rest:        rest,
	})

	normalizeEndpointMeta(s, app, ep)

	// Method params
	for pi, p := range ep.Param {
		normalizeParam(s, app, ep, p.Name, p.Type, pi, "")
	}

	// REST params
	if ep.RestParams != nil {
		for pi, p := range ep.RestParams.UrlParam {
			normalizeParam(s, app, ep, p.Name, p.Type, pi, "path")
		}
		for pi, p := range ep.RestParams.QueryParam {
			normalizeParam(s, app, ep, p.Name, p.Type, pi, "query")
		}
	}

	for i, stmt := range ep.Stmt {
		if err := normalizeStatement(s, app, ep, stmt, i); err != nil {
			return err
		}
	}

	return nil
}

func normalizeEvent(s *Schema, app *sysl.Application, event *sysl.Endpoint) {
	s.Event = append(s.Event, Event{
		AppName:   app.Name.Part,
		EventName: event.Name,
	})

	for pi, p := range event.Param {
		normalizeParam(s, app, event, p.Name, p.Type, pi, "")
	}

	normalizeEventMeta(s, app, event)
}

func normalizeStatement(
	s *Schema,
	app *sysl.Application,
	ep *sysl.Endpoint,
	stmt *sysl.Statement,
	stmtIndex int,
) error {
	statement := Statement{
		AppName:    app.Name.Part,
		EpName:     ep.Name,
		StmtIndex:  stmtIndex,
		StmtParent: StatementParent{},
	}

	if stmt.GetAction() != nil {
		if stmt.GetAction().Action == placeholder {
			return nil
		}
		statement.StmtAction = stmt.GetAction().Action
	}
	if stmt.GetCall() != nil {
		statement.StmtCall = tuple{"appName": app.Name.Part, "epName": ep.Name}
	}
	if stmt.GetCond() != nil {
		statement.StmtCond = tuple{"test": stmt.GetCond().Test}
	}
	if stmt.GetLoop() != nil {
		loop := stmt.GetLoop()
		statement.StmtLoop = tuple{"mode": loop.Mode.String(), "criterion": loop.Criterion}
	}
	if stmt.GetLoopN() != nil {
		statement.StmtLoopN = tuple{"count": stmt.GetLoopN().Count}
	}
	if stmt.GetForeach() != nil {
		statement.StmtForeach = tuple{"coll": stmt.GetForeach().Collection}
	}
	if stmt.GetAlt() != nil {
		statement.StmtAlt = tuple{"choice": ""} // stmt.GetAlt().GetChoice()
	}
	if stmt.GetGroup() != nil {
		statement.StmtGroup = tuple{"title": stmt.GetGroup().Title}
	}
	if stmt.GetRet() != nil && stmt.GetRet().Payload != "" {
		r, err := parseReturnPayload(stmt.GetRet().Payload)
		if err != nil {
			return err
		}
		statement.StmtRet = r
	}
	s.Stmt = append(s.Stmt, statement)

	normalizeStatementMeta(s, app, ep, stmt, stmtIndex)
	return nil
}

func normalizeParam(
	s *Schema,
	app *sysl.Application,
	ep *sysl.Endpoint,
	paramName string,
	paramType *sysl.Type,
	paramIndex int,
	paramLoc string,
) {
	if paramLoc == "" {
		paramLoc = "method"
		if paramType != nil {
			tags := tags(paramType.Attrs)
			if len(tags) > 0 {
				paramLoc = tags[0]
			}
		}
	}

	param := Param{
		AppName:    app.Name.Part,
		EpName:     ep.Name,
		ParamName:  paramName,
		ParamLoc:   paramLoc,
		ParamIndex: paramIndex,
	}
	if paramType == nil {
		param.ParamOpt = false
	} else {
		param.ParamType = parseFieldType(app.Name.Part, paramType)
		param.ParamOpt = paramType.Opt

		normalizeParamMeta(s, app, ep, paramName, paramType)
	}

	s.Param = append(s.Param, param)
}

func normalizeType(s *Schema, app *sysl.Application, typ *sysl.Type, typeName string) {
	s.Type = append(s.Type, Type{
		AppName:       app.Name.Part,
		TypeName:      typeName,
		TypeDocstring: typ.Docstring,
		TypeOpt:       typ.Opt,
	})

	var fields map[string]*sysl.Type
	switch tv := typ.Type.(type) {
	case *sysl.Type_Tuple_:
		fields = tv.Tuple.AttrDefs

	case *sysl.Type_Relation_:
		table := Table{
			AppName:  app.Name.Part,
			TypeName: typeName,
		}
		if typ.GetRelation().PrimaryKey != nil {
			table.Pk = typ.GetRelation().PrimaryKey.AttrName
		}
		s.Table = append(s.Table, table)
		fields = tv.Relation.AttrDefs

	case *sysl.Type_Primitive_, *sysl.Type_Sequence, *sysl.Type_Set, *sysl.Type_TypeRef:
		s.Alias = append(s.Alias, Alias{
			AppName:   app.Name.Part,
			TypeName:  typeName,
			AliasType: parseFieldType(app.Name.Part, typ),
		})

	case *sysl.Type_Enum_:
		e := Enum{
			AppName:   app.Name.Part,
			TypeName:  typeName,
			EnumItems: typ.GetEnum().Items,
		}
		s.Enum = append(s.Enum, e)
	}

	for fieldName, field := range fields {
		normalizeField(s, app, typeName, field, fieldName)
	}

	normalizeTypeMeta(s, app, typ, typeName)
}

func normalizeField(s *Schema, app *sysl.Application, typeName string, field *sysl.Type, fieldName string) {
	fc := FieldConstraint{}
	if field.Constraint != nil {
		for _, c := range field.Constraint {
			if c.Length != nil {
				fc.Length = FieldConstraintLength{
					Min: c.Length.Min,
					Max: c.Length.Max,
				}
			}
			fc.Precision = c.Precision
			fc.Scale = c.Scale
		}
	}

	s.Field = append(s.Field, Field{
		AppName:         app.Name.Part,
		TypeName:        typeName,
		FieldName:       fieldName,
		FieldOpt:        field.Opt,
		FieldType:       parseFieldType(app.Name.Part, field),
		FieldConstraint: fc,
	})

	normalizeFieldMeta(s, app, typeName, field, fieldName)
}

func normalizeView(s *Schema, app *sysl.Application, view *sysl.View, viewName string) {
	s.View = append(s.View, View{
		AppName:  app.Name.Part,
		ViewName: viewName,
		ViewType: parseFieldType(app.Name.Part, view.RetType),
	})

	normalizeViewMeta(s, app, view, viewName)
}

func normalizeAppMeta(s *Schema, app *sysl.Application) {
	tags := tags(app.Attrs)
	for _, tag := range tags {
		s.Tag.App = append(s.Tag.App, AppTag{
			AppName: app.Name.Part,
			AppTag:  tag,
		})
	}

	annos := annos(app.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.App = append(s.Anno.App, AppAnnotation{
			AppName:      app.Name.Part,
			AppAnnoName:  annoName,
			AppAnnoValue: annoValue,
		})
	}

	if len(app.SourceContexts) > 0 {
		s.Src.App = append(s.Src.App, AppContext{
			AppName: app.Name.Part,
			AppSrc:  relmodSourceContext(app.SourceContexts[0]),
			AppSrcs: relmodSourceContexts(app.SourceContexts),
		})
	}
}

func normalizeMixinMeta(s *Schema, app *sysl.Application, mixin *sysl.Application) {
	tags := tags(mixin.Attrs)
	for _, tag := range tags {
		s.Tag.Mixin = append(s.Tag.Mixin, MixinTag{
			AppName:   app.Name.Part,
			MixinName: mixin.Name.Part,
			MixinTag:  tag,
		})
	}

	annos := annos(mixin.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.Mixin = append(s.Anno.Mixin, MixinAnnotation{
			AppName:        app.Name.Part,
			MixinName:      mixin.Name.Part,
			MixinAnnoName:  annoName,
			MixinAnnoValue: annoValue,
		})
	}

	if len(mixin.SourceContexts) > 0 {
		s.Src.Mixin = append(s.Src.Mixin, MixinContext{
			AppName:   app.Name.Part,
			MixinName: mixin.Name.Part,
			MixinSrc:  relmodSourceContext(mixin.SourceContexts[0]),
			MixinSrcs: relmodSourceContexts(mixin.SourceContexts),
		})
	}
}

func normalizeEndpointMeta(s *Schema, app *sysl.Application, ep *sysl.Endpoint) {
	tags := tags(ep.Attrs)
	for _, tag := range tags {
		s.Tag.Ep = append(s.Tag.Ep, EndpointTag{
			AppName: app.Name.Part,
			EpName:  ep.Name,
			EpTag:   tag,
		})
	}

	annos := annos(ep.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.Ep = append(s.Anno.Ep, EndpointAnnotation{
			AppName:     app.Name.Part,
			EpName:      ep.Name,
			EpAnnoName:  annoName,
			EpAnnoValue: annoValue,
		})
	}

	if len(app.SourceContexts) > 0 {
		s.Src.Ep = append(s.Src.Ep, EndpointContext{
			AppName: app.Name.Part,
			EpName:  ep.Name,
			EpSrc:   relmodSourceContext(ep.SourceContexts[0]),
			EpSrcs:  relmodSourceContexts(ep.SourceContexts),
		})
	}
}

func normalizeEventMeta(s *Schema, app *sysl.Application, event *sysl.Endpoint) {
	tags := tags(event.Attrs)
	for _, tag := range tags {
		s.Tag.Event = append(s.Tag.Event, EventTag{
			AppName:   app.Name.Part,
			EventName: event.Name,
			EventTag:  tag,
		})
	}

	annos := annos(event.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.Event = append(s.Anno.Event, EventAnnotation{
			AppName:        app.Name.Part,
			EventName:      event.Name,
			EventAnnoName:  annoName,
			EventAnnoValue: annoValue,
		})
	}

	if len(event.SourceContexts) > 0 {
		s.Src.Event = append(s.Src.Event, EventContext{
			AppName:   app.Name.Part,
			EventName: event.Name,
			EventSrc:  relmodSourceContext(event.SourceContexts[0]),
			EventSrcs: relmodSourceContexts(event.SourceContexts),
		})
	}
}

func normalizeStatementMeta(s *Schema, app *sysl.Application, ep *sysl.Endpoint, stmt *sysl.Statement, stmtIndex int) {
	tags := tags(stmt.Attrs)
	for _, tag := range tags {
		s.Tag.Stmt = append(s.Tag.Stmt, StatementTag{
			AppName:   app.Name.Part,
			EpName:    ep.Name,
			StmtIndex: stmtIndex,
			StmtTag:   tag,
		})
	}

	annos := annos(stmt.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.Stmt = append(s.Anno.Stmt, StatementAnnotation{
			AppName:       app.Name.Part,
			EpName:        ep.Name,
			StmtIndex:     stmtIndex,
			StmtAnnoName:  annoName,
			StmtAnnoValue: annoValue,
		})
	}

	if len(stmt.SourceContexts) > 0 {
		s.Src.Stmt = append(s.Src.Stmt, StatementContext{
			AppName:   app.Name.Part,
			EpName:    ep.Name,
			StmtIndex: stmtIndex,
			StmtSrc:   relmodSourceContext(stmt.SourceContexts[0]),
			StmtSrcs:  relmodSourceContexts(stmt.SourceContexts),
		})
	}
}

func normalizeParamMeta(s *Schema, app *sysl.Application, ep *sysl.Endpoint, paramName string, param *sysl.Type) {
	tags := tags(param.Attrs)
	for _, tag := range tags {
		s.Tag.Param = append(s.Tag.Param, ParamTag{
			AppName:   app.Name.Part,
			EpName:    ep.Name,
			ParamName: paramName,
			ParamTag:  tag,
		})
	}

	annos := annos(param.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.Param = append(s.Anno.Param, ParamAnnotation{
			AppName:        app.Name.Part,
			EpName:         ep.Name,
			ParamName:      paramName,
			ParamAnnoName:  annoName,
			ParamAnnoValue: annoValue,
		})
	}

	if len(param.SourceContexts) > 0 {
		s.Src.Param = append(s.Src.Param, ParamContext{
			AppName:   app.Name.Part,
			EpName:    ep.Name,
			ParamName: paramName,
			ParamSrc:  relmodSourceContext(param.SourceContexts[0]),
			ParamSrcs: relmodSourceContexts(param.SourceContexts),
		})
	}
}

func normalizeTypeMeta(s *Schema, app *sysl.Application, typ *sysl.Type, typeName string) {
	tags := tags(typ.Attrs)
	for _, tag := range tags {
		s.Tag.Type = append(s.Tag.Type, TypeTag{
			AppName:  app.Name.Part,
			TypeName: typeName,
			TypeTag:  tag,
		})
	}

	annos := annos(typ.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.Type = append(s.Anno.Type, TypeAnnotation{
			AppName:       app.Name.Part,
			TypeName:      typeName,
			TypeAnnoName:  annoName,
			TypeAnnoValue: annoValue,
		})
	}

	if len(typ.SourceContexts) > 0 {
		s.Src.Type = append(s.Src.Type, TypeContext{
			AppName:  app.Name.Part,
			TypeName: typeName,
			TypeSrc:  relmodSourceContext(typ.SourceContexts[0]),
			TypeSrcs: relmodSourceContexts(typ.SourceContexts),
		})
	}
}

func normalizeFieldMeta(s *Schema, app *sysl.Application, typeName string, field *sysl.Type, fieldName string) {
	tags := tags(field.Attrs)
	for _, tag := range tags {
		s.Tag.Field = append(s.Tag.Field, FieldTag{
			AppName:   app.Name.Part,
			TypeName:  typeName,
			FieldName: fieldName,
			FieldTag:  tag,
		})
	}

	annos := annos(field.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.Field = append(s.Anno.Field, FieldAnnotation{
			AppName:        app.Name.Part,
			TypeName:       typeName,
			FieldName:      fieldName,
			FieldAnnoName:  annoName,
			FieldAnnoValue: annoValue,
		})
	}

	if len(field.SourceContexts) > 0 {
		s.Src.Field = append(s.Src.Field, FieldContext{
			AppName:   app.Name.Part,
			TypeName:  typeName,
			FieldName: fieldName,
			FieldSrc:  relmodSourceContext(field.SourceContexts[0]),
			FieldSrcs: relmodSourceContexts(field.SourceContexts),
		})
	}
}

func normalizeViewMeta(s *Schema, app *sysl.Application, view *sysl.View, viewName string) {
	tags := tags(view.Attrs)
	for _, tag := range tags {
		s.Tag.View = append(s.Tag.View, ViewTag{
			AppName:  app.Name.Part,
			ViewName: viewName,
			ViewTag:  tag,
		})
	}

	annos := annos(view.Attrs)
	for annoName, annoValue := range annos {
		s.Anno.View = append(s.Anno.View, ViewAnnotation{
			AppName:       app.Name.Part,
			ViewName:      viewName,
			ViewAnnoName:  annoName,
			ViewAnnoValue: annoValue,
		})
	}

	if len(view.SourceContexts) > 0 {
		s.Src.View = append(s.Src.View, ViewContext{
			AppName:  app.Name.Part,
			ViewName: viewName,
			ViewSrc:  relmodSourceContext(view.SourceContexts[0]),
			ViewSrcs: relmodSourceContexts(view.SourceContexts),
		})
	}
}
