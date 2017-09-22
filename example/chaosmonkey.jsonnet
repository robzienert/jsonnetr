{
  chaosMonkey: {
    enabled: true,
    minTimeBetweenKillsInWorkDays: 1,
    exceptions: [
      { region: "*", account: i, detail: "*", stack: "*" },
      for i in ["mgmt", "test", "prod"]
    ],
    regionsAreIndependent: true,
    meanTimeBetweenKillsInWorkDays: 2,
    grouping: "cluster"
  }
}
