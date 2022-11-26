# go-comment-age

Export data about age of comments and functions and run analysers.

## Motivation

Inspired by "Clean Code", Robert C. Martin.
Code comments tend to not be updated when code is updated.
This tool helps to identify and estimate such drift.

TODO: any GenDecl that contains CommentGroup

### Requirements

```bash
go install github.com/nikolaydubina/go-comment-age
```

### Examples

#### [kubernetes](https://github.com/kubernetes/kubernetes)

TODO

### Heuristics

#### Simple Age Difference

Difference of age (time/commit) of any modification of comment from age (time/commit) of any modification of function.

### References

* https://git-scm.com/docs/git-blame
* https://github.com/nishanths/exhaustive
* https://github.com/kubernetes/kubernetes

### Appendix A: Weighted Age

Code changes happen at various rates.
Comments and code can change one line or can change 90% of lines.
It is useful to differentiate between updates.

Implementing this heuristic is more copmlex and requires more tuning.
Thus, it is area of futher research.