(import "default.jsonnet") + (import "delivery-engineering.jsonnet") + {
  name: "orca",
  description: "orchestration service for spinnaker",
  repoSlug: "orca-nflx",
  propertyRolloutConfigId: "abcd1234",
} + (import "chaosmonkey.jsonnet")
