module.exports = {
    preMajor: true,
    issueUrlFormat: "https://github.com/keptn/keptn/issues/{{id}}",
    scripts: {
        postchangelog: "./gh-actions-scripts/post-changelog-actions.sh"
    },
    types: [
        {
            type: "feat",
            section: "Features"
        },
        {
            type: "fix",
            section: "Bug Fixes"
        },
        {
            type: "chore",
            section: "Other"
        },
        {
            type: "docs",
            section: "Docs"
        },
        {
            type: "perf",
            section: "Performance"
        },
        {
            type: "build",
            hidden: true
        },
        {
            type: "ci",
            hidden: true
        },
        {
            type: "refactor",
            section: "Refactoring"
        },
        {
            type: "revert",
            hidden: true
        },
        {
            type: "style",
            hidden: true
        },
        {
            type: "test",
            hidden: true
        }
    ]
};
