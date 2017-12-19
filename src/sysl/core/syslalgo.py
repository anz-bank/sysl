"""sysl algorithms"""
import re
from sysl.core import syslx
from sysl.util import rex

from sysl.proto import sysl_pb2


def enumerate_calls(stmts):  # pylint: disable=too-many-branches
    """Enumerate all calls under stmts, in the form (parent_stmt, call)."""
    for stmt in stmts:
        if stmt.HasField('call'):
            yield (stmt, stmt.call)
        elif stmt.HasField('action'):
            pass
        elif stmt.HasField('cond'):
            for call in enumerate_calls(stmt.cond.stmt):
                yield call
        elif stmt.HasField('loop'):
            for call in enumerate_calls(stmt.loop.stmt):
                yield call
        elif stmt.HasField('loop_n'):
            for call in enumerate_calls(stmt.loop_n.stmt):
                yield call
        elif stmt.HasField('foreach'):
            for call in enumerate_calls(stmt.foreach.stmt):
                yield call
        elif stmt.HasField('group'):
            for call in enumerate_calls(stmt.group.stmt):
                yield call
        elif stmt.HasField('alt'):
            for choice in stmt.alt.choice:
                for call in enumerate_calls(choice.stmt):
                    yield call
        elif stmt.HasField('ret'):
            pass
        else:
            raise Exception('No statement!', stmt.WhichOneof('stmt'))


# TODO: Require consistency along different branches.
def return_payload(stmts):
    # pylint: disable=too-many-branches,too-many-return-statements
    """Compute payload returned by stmts."""
    for stmt in stmts:
        if stmt.HasField('call'):
            pass
        elif stmt.HasField('action'):
            pass
        elif stmt.HasField('cond'):
            payload = return_payload(stmt.cond.stmt)
            if payload:
                return payload
        elif stmt.HasField('loop'):
            payload = return_payload(stmt.loop.stmt)
            if payload:
                return payload
        elif stmt.HasField('loop_n'):
            payload = return_payload(stmt.loop_n.stmt)
            if payload:
                return payload
        elif stmt.HasField('foreach'):
            payload = return_payload(stmt.foreach.stmt)
            if payload:
                return payload
        elif stmt.HasField('group'):
            payload = return_payload(stmt.group.stmt)
            if payload:
                return payload
        elif stmt.HasField('alt'):
            for choice in stmt.alt.choice:
                payload = return_payload(choice.stmt)
                if payload:
                    return payload
        elif stmt.HasField('ret'):
            return stmt.ret.payload
        else:
            raise Exception('No statement!', stmt.WhichOneof('stmt'))


def yield_ret_params(payload):
    if payload is not None:
        for param_pair in rex.split(r',(?![^{]*\})', payload):
            ptname = param_pair

            if param_pair.count('<:') == 1:
                (pname, ptname) = rex.split(r'\s*<:\s*', param_pair)

            if ptname.upper() not in (
                    set(sysl_pb2.Type.Primitive.keys()) - {'NO_Primitive'}):
                m = rex.match(r'set\s+of\s+(.+)$', ptname)
                if m:
                    ptname = m.group(1)

                m = rex.match(r'one\s+of\s*{(.+)}$', ptname)
                if m:
                    for ptn in rex.split(r'\s*,\s*', m.group(1)):
                        yield ptn
                else:
                    yield ptname
