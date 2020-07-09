# Data Flow Analysis

Generates a data flow analysis for a field of `model.sysl`.

```bash
$ arrai run --data=Source.Customer.name data_flow.arrai
Source.Customer.name <- Source.Select <- A.Fetch
A.FetchResponse.user_name <- D.Fetch
# TODO: C.FetchResponse.
D.FetchResponse.userName <- Client.ProfileScreen
Client.Profile.displayName
```
