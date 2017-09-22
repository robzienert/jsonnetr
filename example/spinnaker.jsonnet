(import "default.jsonnet") + (import "delivery-engineering.jsonnet") + {
  name: "orca",
  description: "orchestration service for spinnaker",
  repoSlug: "orca-nflx",
  propertyRolloutConfigId: "abcd1234",
} + (import "chaosmonkey.jsonnet") + (import "https://gist.githubusercontent.com/robzienert/485938aaafcd42768923b7e02cb811a7/raw/322f11adae817a48c9b28e1c3c666909308a320e/pipelines.jsonnet")
