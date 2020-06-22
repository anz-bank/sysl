# Sysl Bank Example

This demo provides an example of a very simple banking system in which a `Customer` of a `Branch` can make a `Transaction` to move money in or out of an `Account`. The bank maintains records of these entities, and offers APIs to read and write customers and accounts.

## Data model diagram

```bash
sysl data --output data.png --project "Bank :: Data Views" ./project.sysl
```

## Integration diagram

```bash
sysl ints --output ints.png --project "Bank :: Integrations" ./project.sysl
```

## Sequence diagram

For the ATM withdrawal endpoint:

```bash
sysl sd --output sd.png --endpoint "ATM <- Withdraw" ./project.sysl
```

## Database script

```bash
sysl generate-db-scripts --app-names BankModel ./bank.sysl
```
