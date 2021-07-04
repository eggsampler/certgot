# `certgot/log`

---

This package aims to implement a simple, semi-structured levelled logger while not importing many (ideally: not 
importing any) external packages.

It is not performant and probably needs to be benchmarked and cleaned up. That said, it's being used in a cli
application and won't be logging fast or extensively, so this is not really a priority for now.

TODO:

- Implement logging to file
  - This would mean logging to a buffer in memory until specifying a file somehow