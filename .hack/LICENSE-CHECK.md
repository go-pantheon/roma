# License Checking Tool User Guide

This project implements a comprehensive open-source dependency license checking mechanism to ensure all dependencies comply with the project's license policy.

## License Checking Process

1. **GitHub Actions Automatic Check**: GitHub Actions automatically runs license checking on every push and PR creation
2. **Local pre-commit Check**: The pre-commit hook automatically checks licenses before committing code
3. **Manual Check**: You can manually run the script to check licenses at any time

## License Whitelist

Currently allowed license types:
- MIT
- Apache-2.0
- BSD-2-Clause
- BSD-3-Clause
- ISC
- MPL-2.0

To modify the whitelist, edit the `ALLOWED_LICENSES` array in the `./.hack/licenses-check` script.

## Available Commands

### Check License Compliance

```bash
./.hack/licenses-check
```

This command checks all dependency licenses and returns a non-zero exit code with a list of problematic dependencies if licenses not in the whitelist are found.

### Generate License Report

```bash
./.hack/licenses-generate
```

This command generates the following files:
- `DEPENDENCIES.md`: Markdown report of dependency licenses
- `license-check/third-party-licenses.<os>.md`: License reports for different operating systems
- `licenses/` directory: Contains copies of all dependency licenses
- `license-check/license-dependencies.csv`: Standard CSV format license report

## Setting up pre-commit

License checking has been added to the pre-commit configuration. To enable it, run:

```bash
pre-commit install
```

## Handling License Issues

If a dependency uses a disallowed license, you can:

1. Find an alternative dependency with an allowed license
2. Get an exception approval and add the dependency to the exception list
3. Update the license whitelist after evaluating the risk

## Common Warnings and Solutions

### Empty Module Version Warning

You may see warnings like:

```
W0414 11:47:02.772394   46191 library.go:276] module github.com/go-pantheon/roma has empty version, defaults to HEAD. The license URL may be incorrect. Please verify!
```

This warning typically appears when:

1. The project hasn't released an official version (no Git Tag)
2. Go's replace directive is used to point to a local directory
3. Go workspace mode is used for multi-module project development

**Solutions**:
- These warnings don't affect the correctness of license checking, only the generated license URL might be imprecise
- The script automatically filters these warnings, logging them to the `license-check-warnings.log` file
- For more precise URLs, create release tags (git tag) for the project, then rerun the check
- **Automatically ignore the project itself**: The script automatically identifies the project's own packages (via module name in go.mod) and ignores related warnings in the report
- If you want to completely exclude the project's own packages from checking, add the project module name to `./.hack/license-exceptions.txt`
