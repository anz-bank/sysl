/** The kinds of Sysl element that an annotation can be applied to. */
export enum ElementKind {
    App = "app",
    Type = "type",
    Field = "field",
    Endpoint = "ep",
    Statement = "stmt",
}

export namespace ElementKind {
    /**
     * Converts an {@link ElementKind} to a numerical form. Applications are 1, types and endpoints are 2, fields and
     * statements are 3.
     * @param kind The {@link ElementKind} for which to return the numeric form.
     * @returns The numerical form for the specified kind.
     */
    export function toNumeric(kind: ElementKind): number {
        if (kind == ElementKind.App) return 1;
        if (kind == ElementKind.Type || kind == ElementKind.Endpoint) return 2;
        return 3;
    }

    const numericToDataKind: readonly ElementKind[] = [ElementKind.App, ElementKind.Type, ElementKind.Field];
    const numericToBehavKind: readonly ElementKind[] = [ElementKind.App, ElementKind.Endpoint, ElementKind.Statement];

    /**
     * Converts the numerical form of an {@link ElementKind} to their corresponding data-related instance: application
     * for 1, type/endpoint for 2 and field/statement for 3.
     * @param numericKind The numerical form of an {@link ElementKind}.
     * @returns The data-related {@link ElementKind} for the corresponding numerical form.
     * @throws `Error` if an invalid numerical form is specified.
     */
    export function fromNumeric(numericKind: number, behavior?: boolean): ElementKind {
        const kind = (behavior ? numericToBehavKind : numericToDataKind)[numericKind - 1];
        if (!kind) throw new Error("Invalid depth. Must be between 1 and 3.");
        return kind;
    }
}
