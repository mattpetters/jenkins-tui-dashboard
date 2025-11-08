# TODO: Fix GitHub Repo Path

## Current Issue

GitHub branch fetch returns null because repo path is incorrect.

## What We Tried

```bash
# Tried: intuit/identity-manage/account
curl "https://github.intuit.com/api/v3/repos/intuit/identity-manage/account/pulls/3934"
# Result: 404 or null

# Tried: intuit/identity-manage  
curl "https://github.intuit.com/api/v3/repos/intuit/identity-manage/pulls/3934"
# Result: null
```

## What We Know

- PR-3934 web URL: `https://github.intuit.com/identity-manage/account/pull/3934`
- This URL works in browser
- But API path structure is different

## Possible Solutions

### Option 1: Find Correct Repo Name
Search GitHub or ask someone for the exact full_name of the repo.

### Option 2: Use GitHub MCP
Instead of direct API calls, use the Intuit GitHub MCP which already works.

### Option 3: Manual Entry
Have user optionally enter Git branch when adding PR (fallback).

## For Now

- Selection highlighting improved ✅
- URLs fixed (Enter/p both work) ✅
- Git branch field exists in model ✅
- Will show "master" from Jenkins or fallback to "PR-XXXX" ✅

GitHub branch auto-fetch can be fixed once we find correct repo path.

