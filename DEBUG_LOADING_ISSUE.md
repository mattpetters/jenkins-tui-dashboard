# Debugging "Loading..." Issue

## What to Check

Run `./run.sh` and add PR-3934 again.

### Look at the Status Message

After the build is fetched, the status bar will show:
```
✓ PR-3934: failure (Stage: XXXX, Job: YYYY)
```

**Check what Stage and Job show:**

1. **If they're empty strings** → Jenkins API isn't returning stages data
2. **If they show actual values** → The tile rendering has a bug
3. **If status never updates** → The fetch command isn't completing

### Jenkins API Debugging

You can manually test the Jenkins API:

```bash
curl -u "mpetters:113e76ce0e6775a1dca5a29d857e025888" \
  "https://build.intuit.com/identity/job/identity-manage/job/account/job/account-eks/job/PR-3934/lastBuild/api/json" \
  | jq '.stages'
```

This will show if Jenkins is returning the stages array.

### GitHub API Debugging

Test if GitHub returns the branch name:

```bash
curl -H "Authorization: Bearer github_pat_11AAASLLY09..." \
  "https://github.intuit.com/api/v3/repos/intuit/identity-manage/account/pulls/3934" \
  | jq '.head.ref'
```

Should return: "IDLMP-2038-aggregate"

## Next Steps

Based on the debug output, we'll know:
- If Jenkins returns stages → Fix the parser
- If Jenkins doesn't return stages → Use different API endpoint
- If GitHub fetch works → Should see branch name
- If GitHub fetch fails → Will see "PR-3934" as fallback

