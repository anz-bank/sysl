# Integration Diagram

Generates the PlantUML for an integration diagram of `model.sysl`.

```bash
$ arrai run integration_diagram.arrai
@startuml
[A] --> [Source]
[B] --> [Source]
[C] --> [A]
[C] --> [B]
[Client] --> [D]
[D] --> [A]
[D] --> [C]
@enduml
```

![Integration diagram](integration_diagram.svg)
