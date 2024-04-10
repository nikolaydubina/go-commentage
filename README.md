# go-commentage

[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/go-commentage)](https://goreportcard.com/report/github.com/nikolaydubina/go-commentage)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nikolaydubina/go-commentage/badge)](https://securityscorecards.dev/viewer/?uri=github.com/nikolaydubina/go-commentage)

How far behind are comments compared to code? Are they being updated?
Inspired by "Clean Code" by Robert C. Martin, this `go vet` compatible tool analyses AST and `git` and collects details on how far comments drift from code they describe.

You need to have `git` with version `>=2.37`.

```bash
go install github.com/nikolaydubina/go-commentage@latest
```

```bash
go-commentage ./...
```

```txt
kubernetes/pkg/util/ipset/ipset.go:283:1: "CreateSet": doc_last_updated_behind_days(1336.83)
kubernetes/pkg/util/ipset/ipset.go:296:1: "createSet": doc_last_updated_behind_days(1603.17)
kubernetes/pkg/util/ipset/ipset.go:320:1: "AddEntry": doc_last_updated_behind_days(1578.10)
kubernetes/pkg/util/ipset/ipset.go:332:1: "DelEntry": doc_last_updated_behind_days(1578.10)
kubernetes/pkg/util/ipset/ipset.go:340:1: "TestEntry": doc_last_updated_behind_days(450.07)
kubernetes/pkg/util/ipset/ipset.go:356:1: "FlushSet": doc_last_updated_behind_days(0.00)
kubernetes/pkg/util/ipset/ipset.go:364:1: "DestroySet": doc_last_updated_behind_days(73.85)
kubernetes/pkg/util/ipset/ipset.go:372:1: "DestroyAllSets": doc_last_updated_behind_days(0.00)
kubernetes/pkg/util/ipset/ipset.go:380:1: "ListSets": doc_last_updated_behind_days(0.00)
kubernetes/pkg/util/ipset/ipset.go:389:1: "ListEntries": doc_last_updated_behind_days(0.00)
```

### Filtering

To narrow down output, filters can be specified.
This allows interactive exploration.
This also allows usage of this tool in CI, since when no diagnostic is printed then status code is 0.

```bash
$ go-commentage -min-days-behind 100 ./...
$ go-commentage -commit -min-days-behind 100 ./...
$ go-commentage -min-days-behind 10 -commit -min-days-behind 100 ./...
$ echo $?
0
```

### Heuristics

#### Simple Age Difference

Measure of how far away in terms of days or commits last update of function body as compared to last update to associated doc comment.

#### Weighted Age Difference

> [!Warning]  
> Work in Progress

Code changes happen at various rates.
Comments and code can change one line or can change 90% of lines.
It is useful to differentiate between updates.

### References

* https://git-scm.com/docs/git-blame
* https://git-scm.com/docs/git-rev-list
* https://github.com/nishanths/exhaustive
* https://github.com/kubernetes/kubernetes

### Appendix A: Full Output

This can be useful for debugging or exporting for post processing and further data visualization.

```bash
go-commentage -verbose -time -commit ./...
```

```txt
kubernetes/pkg/util/pod/pod.go:34:1: "PatchPodStatus": last_updated_at(2022-08-02T13:58:08+08:00) doc_last_updated_at(2020-02-27T06:05:33+08:00) doc_last_updated_behind_days(887.33)last_commit(04fcbd721cd3) doc_last_commit(b2528654797e) doc_last_commit_behind(8786)
kubernetes/pkg/util/pod/pod.go:74:1: "ReplaceOrAppendPodCondition": last_updated_at(2022-11-07T18:57:56+08:00) doc_last_updated_at(2022-11-07T18:57:56+08:00) doc_last_updated_behind_days(0.00)last_commit(4e732e20d05e) doc_last_commit(4e732e20d05e) doc_last_commit_behind(0)
kubernetes/pkg/util/procfs/procfs_fake.go:28:1: "GetFullContainerName": last_updated_at(2015-11-14T07:47:25+08:00) doc_last_updated_at(2017-04-04T14:16:34+08:00) doc_last_updated_behind_days(-507.27)last_commit(fb576f30c838) doc_last_commit(932ece5cfd0f) doc_last_commit_behind(-10936)
kubernetes/pkg/util/procfs/procfs_unsupported.go:34:1: "GetFullContainerName": last_updated_at(2016-08-18T23:01:03+08:00) doc_last_updated_at(2016-08-17T07:34:14+08:00) doc_last_updated_behind_days(1.64)last_commit(5eef6b8d91a2) doc_last_commit(a2824bb7a337) doc_last_commit_behind(58)
kubernetes/pkg/util/procfs/procfs_unsupported.go:40:1: "PKill": last_updated_at(2016-08-18T23:01:03+08:00) doc_last_updated_at(2016-08-18T23:01:03+08:00) doc_last_updated_behind_days(0.00)last_commit(5eef6b8d91a2) doc_last_commit(5eef6b8d91a2) doc_last_commit_behind(0)
kubernetes/pkg/util/procfs/procfs_unsupported.go:46:1: "PidOf": last_updated_at(2016-08-18T23:01:03+08:00) doc_last_updated_at(2016-08-18T23:01:03+08:00) doc_last_updated_behind_days(0.00)last_commit(5eef6b8d91a2) doc_last_commit(5eef6b8d91a2) doc_last_commit_behind(0)
kubernetes/pkg/util/removeall/removeall.go:35:1: "RemoveAllOneFilesystemCommon": last_updated_at(2021-06-04T06:38:37+08:00) doc_last_updated_at(2021-06-04T06:38:37+08:00) doc_last_updated_behind_days(0.00)last_commit(484eb0182224) doc_last_commit(484eb0182224) doc_last_commit_behind(0)
kubernetes/pkg/util/removeall/removeall.go:115:1: "RemoveAllOneFilesystem": last_updated_at(2021-06-04T06:38:37+08:00) doc_last_updated_at(2021-06-16T00:40:17+08:00) doc_last_updated_behind_days(-11.75)last_commit(484eb0182224) doc_last_commit(01bb0f86b02b) doc_last_commit_behind(-1)
kubernetes/pkg/util/removeall/removeall.go:126:1: "RemoveDirsOneFilesystem": last_updated_at(2021-06-04T06:38:37+08:00) doc_last_updated_at(2021-06-16T00:40:17+08:00) doc_last_updated_behind_days(-11.75)last_commit(484eb0182224) doc_last_commit(01bb0f86b02b) doc_last_commit_behind(-1)
kubernetes/pkg/util/rlimit/rlimit_unsupported.go:27:1: "SetNumFiles": last_updated_at(2020-02-25T13:58:28+08:00) doc_last_updated_at(2020-02-25T13:58:28+08:00) doc_last_updated_behind_days(0.00)last_commit(4936cd476bf3) doc_last_commit(4936cd476bf3) doc_last_commit_behind(0)
kubernetes/pkg/util/slice/slice.go:26:1: "CopyStrings": last_updated_at(2017-06-23T11:41:18+08:00) doc_last_updated_at(2015-01-23T06:12:37+08:00) doc_last_updated_behind_days(882.23)last_commit(f98bc7d45435) doc_last_commit(f7e3cb12a6e7) doc_last_commit_behind(19409)
kubernetes/pkg/util/slice/slice.go:37:1: "SortStrings": last_updated_at(2015-01-23T06:12:37+08:00) doc_last_updated_at(2015-01-23T06:12:37+08:00) doc_last_updated_behind_days(0.00)last_commit(f7e3cb12a6e7) doc_last_commit(f7e3cb12a6e7) doc_last_commit_behind(0)
kubernetes/pkg/util/slice/slice.go:44:1: "ContainsString": last_updated_at(2017-04-07T08:14:16+08:00) doc_last_updated_at(2017-04-07T08:14:16+08:00) doc_last_updated_behind_days(0.00)last_commit(151770c8fde9) doc_last_commit(151770c8fde9) doc_last_commit_behind(0)
kubernetes/pkg/util/slice/slice.go:58:1: "RemoveString": last_updated_at(2017-11-23T23:00:35+08:00) doc_last_updated_at(2017-11-23T23:00:35+08:00) doc_last_updated_behind_days(0.00)last_commit(e1312f2c00ed) doc_last_commit(e1312f2c00ed) doc_last_commit_behind(0)
kubernetes/pkg/util/tail/tail.go:38:1: "ReadAtMost": last_updated_at(2022-10-20T15:13:28+08:00) doc_last_updated_at(2016-12-08T04:56:06+08:00) doc_last_updated_behind_days(2142.43)last_commit(cc90e819bce9) doc_last_commit(2bb2604f0b0d) doc_last_commit_behind(29079)
kubernetes/pkg/util/tail/tail.go:68:1: "FindTailLineStartIndex": last_updated_at(2018-02-11T11:02:23+08:00) doc_last_updated_at(2016-12-08T04:56:06+08:00) doc_last_updated_behind_days(430.25)last_commit(7cfb94cbc576) doc_last_commit(2bb2604f0b0d) doc_last_commit_behind(8120)
kubernetes/pkg/util/tolerations/tolerations.go:27:1: "VerifyAgainstWhitelist": last_updated_at(2019-08-21T09:21:57+08:00) doc_last_updated_at(2017-02-28T02:34:46+08:00) doc_last_updated_behind_days(904.28)last_commit(5a50b3f4a2a2) doc_last_commit(af5379485411) doc_last_commit_behind(15293)
```
