import { spawn, SpawnOptions } from "child_process";

/** Spawns a child process and returns a promise that resolves to the buffer content of stdout. */
export function spawnBuffer(
    command: string,
    args: ReadonlyArray<string> = [],
    options: SpawnOptions & { input?: any } = {}
): Promise<Buffer> {
    return new Promise<Buffer>((resolve, reject) => {
        const process = spawn(command, args, options);
        if (options.input) {
            options.input && process.stdin?.write(options.input);
        }
        process.stdin?.end();

        let chunks: any[] = [];
        let result: Buffer;
        process.stdout?.on("data", (data: Buffer) => chunks.push(data));
        process.stdout?.on("end", () => (result = Buffer.concat(chunks)));
        let err = "";
        process.stderr?.on("data", data => {
            err += data;
        });
        process.on("error", err => reject(`spawn error: ${err}`));
        process.on("close", code => {
            code === 0 ? resolve(result) : reject(`spawn exited with code ${code}:\n${err}`);
        });
    });
}

export async function spawnSysl(args: string[], options: SpawnOptions & { input?: any } = {}): Promise<Buffer> {
    const syslPath = process.env["SYSL_PATH"] ?? "sysl";
    return await spawnBuffer(syslPath, args, options);
}
