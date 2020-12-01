# Testing

Sysl provides `test-rig` tool for quick bootstrapping RESTful service along with database component, once you used `codegen` tool and implemented your business logic.

```text
sysl test-rig --template variables.json --output-dir out
```

This tool generates boilerplate for single (yet) service based on sysl app definitions in sysl module(-s) and substitutions provided in `--template` parameter. The result is ready-to-use docker-compose script, that includes containers needed to run your business logic as REST http service.

## Design

For each app found in sysl module, `test-rig` generates Go program out of pre-baked template - basically, `main()` func with minimum boilerplate. Next, `test-rig` generates Docker script that builds this Go program and packs it to lightweight container. If your app has [~DB] attribute and `table` blocks, `test-rig` generates sidecar postgres container with baked-in schema initialization script, and sets up connection.

## Howto

### Sysl routines

1. Run `sysl codegen` with proper transform and grammar to get RESTful service stubs
2. Run `sysl generate-db-scripts` if you have apps that require databases
3. Run `sysl codegen` for downstream apps if you have any
4. And then you can run `sysl test-rig`

### User routines

1. `sysl test-rig` requires you to provide two factory functions. One to return `GenCallback` inteface implementation:

   ```go
   // non-db version
   type factoryFuncType = func() GenCallback
   ```

   ```go
   // db version
   type factoryFuncType = func(db *sql.DB) GenCallback
   ```

   Another to return `ServiceInterface` implementation:

   ```go
   type factoryFuncType = func() ServiceInterface
   ```

2. `test-rig` needs json file that contains things needed for `main()` templating. This json file contains number of blocks, each representing sysl app you implemented business logic for. Example:

   ```json
   {
       "dbfront": {  <- sysl app name
           "name": "dbfront", <- generated package name for your app; this also will be your container name
           "import": "< golang import path to generated code >",
           "port": "8080", <- port for http server
           "impl": {
               "name": "dbfront_impl", <- package name for your implementation
               "import": "< golang import path to your implementation >",
               "interface_factory": "GetServiceInterfaceImplementation", <- factory function name exported from your implementation package
               "callback_factory": "GetCallback" <- factory function name exported from your implementation package
           }
       }
   }
   ```

## Example project

https://github.com/anz-bank/test-rig-example
