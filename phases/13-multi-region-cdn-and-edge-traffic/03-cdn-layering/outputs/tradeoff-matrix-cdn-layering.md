---
lesson: 03-cdn-layering
---

| Layer | Main benefit | Main cost | Best fit |
|------|--------------|-----------|----------|
| Browser | zero network for repeats | weak central control | immutable or user-local content |
| Edge POP | low user latency | POP memory and key complexity | global hot content |
| Shield | origin protection | extra hop and ops surface | many POPs, few origins |
| Origin-adjacent | backend smoothing | least user-latency benefit | expensive origin fetches |
