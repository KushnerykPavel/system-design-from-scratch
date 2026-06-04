# Interview Card — SLIs, SLOs, and Error Budgets

## Fast answer shape

1. Name the user journey.
2. Define the SLI numerator and denominator.
3. Set a realistic target and time window.
4. Explain the error budget and burn policy.
5. Describe what decisions change when the budget burns.

## Good phrases

- "I want the SLI at the user-facing boundary, not only at the host boundary."
- "A target without an action policy is just a reporting number."
- "I would separate read and write objectives if their value or failure shape differs."

## Weak answer smells

- uptime as the only objective
- no numerator or denominator
- no latency consideration for a latency-sensitive path
- no explanation of what budget burn changes operationally
