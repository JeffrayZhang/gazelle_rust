# Notes

## Resolve
> Objective: Fill the "deps" attribute of a rule with the list of its dependencies
* check for `gazelle:resolve` "overrides" using `resolve.FindRuleWithOverride`
* check for third-party dependencies with `cfg.FindThirdPartyDependency`
* check if it's a standard library dependency (ignore if so)
* return an error if a dependency is not found (and `cfg.ValidateImportStatements` is true; else ignore)
* convert `resolve.ImportSpec` to `[]resolve.FindResult` by looking up with `ix.FindRulesByImportWithConfig`
  * filter out self-dependencies that would cause an infinite loop (i.e. `//foo:bar` cannot depend on `//foo:bar`)
  * return an error if multiple results are found (i.e. the program is unsure which target to use as dependency)
* set deps attribute with `r.SetAttr("deps", ...)`

## GenerateRules
> Objective: Given a directory, generate a list of rules that correspond to the files in that directory. Optionally, generate metadata for each of those rules.

## Imports
> Objective: Assign one or more keys (a.k.a. ImportSpecs) to a given rule/file, so that we can find it later with `ix.FindRulesByImportWithConfig`
