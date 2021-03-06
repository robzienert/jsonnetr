# example

This is a rather contrived example of using `jsonnetr` to build a Spinnaker 
application json definition. The root file is `spinnaker.jsonnet`, which has
a few local-file imports, as well as a remote HTTP import. Running `jsonnetr`
will yield the following output.

The remote HTTP import just includes a `pipelines` key with an empty list.

```
$ jsonnetr spinnaker.jsonnet
```

```json
{
   "chaosMonkey": {
      "enabled": true,
      "exceptions": [
         {
            "account": "mgmt",
            "detail": "*",
            "region": "*",
            "stack": "*"
         },
         {
            "account": "test",
            "detail": "*",
            "region": "*",
            "stack": "*"
         },
         {
            "account": "prod",
            "detail": "*",
            "region": "*",
            "stack": "*"
         }
      ],
      "grouping": "cluster",
      "meanTimeBetweenKillsInWorkDays": 2,
      "minTimeBetweenKillsInWorkDays": 1,
      "regionsAreIndependent": true
   },
   "dataSources": {
      "disabled": [ ],
      "enabled": [
         "analytics"
      ]
   },
   "description": "orchestration service for spinnaker",
   "email": "example@example.com",
   "enableRestartRunningExecutions": false,
   "group": "Spinnaker",
   "instancePort": 7001,
   "name": "orca",
   "owner": "Delivery Engineering",
   "pdApiKey": "abcd1234",
   "pipelines": [ ],
   "propertyRolloutConfigId": "abcd1234",
   "repoProjectKey": "SPKR",
   "repoSlug": "orca-nflx",
   "repoType": "stash",
   "requiredGroupMembership": [ ],
   "type": "Web Service"
}
```

