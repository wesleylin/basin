## development

go test ./...

go test ./orderedmap

go test -bench=. -benchmem

go test -bench=. -benchmem > bench_results.txt

git tag -a v0.1.0 -m "xyz"

git push origin v0.1.0
