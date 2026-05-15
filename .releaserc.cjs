module.exports = {
  extends: ["semantic-release-config-gitmoji"],
  branches: [
    "main",
    { name: "next", prerelease: "rc" },
    { name: "beta", prerelease: true },
    { name: "alpha", prerelease: true }
  ]
};
