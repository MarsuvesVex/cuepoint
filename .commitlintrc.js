const pkg = require('./package.json')
/**/
/* // Check if the user has configured the package to use conventional commits. */
const isConventional = pkg.config ? pkg.config['cz-emoji']?.conventional : false

// Regex for default and conventional commits.
const RE_DEFAULT_COMMIT = /^(?::.*:|(?:\u00a9|\u00ae|[\u2000-\u3300]|\ud83c[\ud000-\udfff]|\ud83d[\ud000-\udfff]|\ud83e[\ud000-\udfff]))\s(?<emoji>\((?<scope>.*)\)\s)?.*$/gm
const RE_CONVENTIONAL_COMMIT = /^^(?<type>\w+)(?:\((?<scope>\w+)\))?\s(?<emoji>:.*:|(?:\u00a9|\u00ae|[\u2000-\u3300]|\ud83c[\ud000-\udfff]|\ud83d[\ud000-\udfff]|\ud83e[\ud000-\udfff]))\s.*$/gm

/* module.exports = { */
/*   rules: { */
/*     'cz-emoji': [2, 'always'] */
/*   }, */

/*   plugins: [ */
/*     { */
/*       rules: { */
/*         'cz-emoji': ({ raw }) => { */
/*           const isValid = isConventional */
/*             ? RE_CONVENTIONAL_COMMIT.test(raw) */
/*             : RE_DEFAULT_COMMIT.test(raw) */
/**/
/*           const message = isConventional */
/*             ? `Your commit message should follow conventional commit format.` */
/*             : `Your commit message should be: <emoji> (<scope>)?: <subject>` */
/**/
/*           return [isValid, message] */
/*         } */
/*       } */
/*     } */
/*   ] */
/* } */
// .commitlintrc.js
/* const fs = require('fs'); */
/* const path = require('path'); */
/* const packages = fs.readdirSync(path.resolve(__dirname, 'packages')); */
/**/
/* module.exports = { */
/*   rules: { */
/*     'scope-enum': [2, 'always', packages] */
/*   }, */
/*   prompt: { */
/*     scopes: [...packages] */
/*   } */
/* }; */
/**/

const { defineConfig } = require('cz-git')

module.exports = defineConfig({
	extends: ['@commitlint/config-conventional'],
  plugins: [
    {
      rules: {
        'cz-emoji': ({ raw }) => {
          const isValid = isConventional
            ? RE_CONVENTIONAL_COMMIT.test(raw)
            : RE_DEFAULT_COMMIT.test(raw)

          const message = isConventional
            ? `Your commit message should follow conventional commit format.`
            : `Your commit message should be: <emoji> (<scope>)?: <subject>`

          return [isValid, message]
        }
      }
    }
  ],
    prompt: {
        alias: { fd: 'docs: fix typos' },
        messages: {
            type: 'Select the type of change that you\'re committing:',
            scope: 'Denote the SCOPE of this change (optional):',
            customScope: 'Denote the SCOPE of this change:',
            subject: 'Write a SHORT, IMPERATIVE tense description of the change:\n',
            body: 'Provide a LONGER description of the change (optional). Use "|" to break new line:\n',
            breaking: 'List any BREAKING CHANGES (optional). Use "|" to break new line:\n',
            footerPrefixesSelect: 'Select the ISSUES type of changeList by this change (optional):',
            customFooterPrefix: 'Input ISSUES prefix:',
            footer: 'List any ISSUES by this change. E.g.: #31, #34:\n',
            generatingByAI: 'Generating your AI commit subject...',
            generatedSelectByAI: 'Select suitable subject by AI generated:',
            confirmCommit: 'Are you sure you want to proceed with the commit above?',
        },
        types: [
            { value: 'feat', name: 'feat:     ✨  A new feature', emoji: ':sparkles:' },
            { value: 'fix', name: 'fix:      🐛  A bug fix', emoji: ':bug:' },
            { value: 'docs', name: 'docs:     📝  Documentation only changes', emoji: ':memo:' },
            { value: 'style', name: 'style:    💄  Changes that do not affect the meaning of the code', emoji: ':lipstick:' },
            { value: 'refactor', name: 'refactor: ♻️   A code change that neither fixes a bug nor adds a feature', emoji: ':recycle:' },
            { value: 'perf', name: 'perf:     ⚡️  A code change that improves performance', emoji: ':zap:' },
            { value: 'test', name: 'test:     ✅  Adding missing tests or correcting existing tests', emoji: ':white_check_mark:' },
            { value: 'build', name: 'build:    📦️   Changes that affect the build system or external dependencies', emoji: ':package:' },
            { value: 'ci', name: 'ci:       🎡  Changes to our CI configuration files and scripts', emoji: ':ferris_wheel:' },
            { value: 'chore', name: 'chore:    🔨  Other changes that don\'t modify src or test files', emoji: ':hammer:' },
            { value: 'revert', name: 'revert:   ⏪️  Reverts a previous commit', emoji: ':rewind:' },
        ],
        useEmoji: true,
        emojiAlign: 'center',
        useAI: false,
        aiNumber: 1,
        themeColorCode: '',
      scopes: [
        "api",
        "bot",
        "worker",
        "cli",
        "config",
        "database",
        "events",
        "ffmpeg",
        "stream",
        "infra",
        "docs",
        "repo"
      ],
        allowCustomScopes: true,
        allowEmptyScopes: true,
        customScopesAlign: 'bottom',
        customScopesAlias: 'custom',
        emptyScopesAlias: 'empty',
        upperCaseSubject: null,
        markBreakingChangeMode: false,
        allowBreakingChanges: ['feat', 'fix'],
        breaklineNumber: 100,
        breaklineChar: '|',
        skipQuestions: [],
        issuePrefixes: [{ value: 'closed', name: 'closed:   ISSUES has been processed' }],
        customIssuePrefixAlign: 'top',
        emptyIssuePrefixAlias: 'skip',
        customIssuePrefixAlias: 'custom',
        allowCustomIssuePrefix: true,
        allowEmptyIssuePrefix: true,
        confirmColorize: true,
        scopeOverrides: undefined,
        defaultBody: '',
        defaultIssues: '',
        defaultScope: '',
        defaultSubject: '',
    },
})
