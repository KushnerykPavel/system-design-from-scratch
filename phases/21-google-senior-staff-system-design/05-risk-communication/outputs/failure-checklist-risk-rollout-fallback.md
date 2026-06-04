# Failure Checklist - Risk, Rollout, and Fallback

## Ask

- what fails first?
- what detects it fastest?
- how small can the canary be?
- what guarantee must survive degradation?

## During Rollout

- compare old path and new path
- pause on clear health gates
- widen only after stable canary behavior

## If Things Go Wrong

- know whether to rollback or fallback
- preserve correctness before convenience
- record why the rollout paused or retreated
