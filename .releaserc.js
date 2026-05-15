module.exports = {
  extends: ["semantic-release-config-gitmoji"],

  branches: ["main",
      { name: "next", channel: "next", prerelease: "rc" },
    { name: "beta", prerelease: true },
    { name: "alpha", prerelease: true }],

  tagFormat: "v${version}",

  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",

    [
      "@semantic-release/changelog",
      {
        changelogFile: "CHANGELOG.md",
      },
    ],

    [
      "@semantic-release/git",
      {
        assets: ["CHANGELOG.md"],
        message:
          "🔖 release: v${nextRelease.version} [skip ci]\n\n${nextRelease.notes}",
      },
    ],

    "@semantic-release/github",
  ],
};
