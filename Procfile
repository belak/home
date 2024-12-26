# Web runs the server which handles both the dynamic portion *and* serving the
# assets.
web: reflex --decoration=none -s -r '\.(go|css|js)$' -- go run ./cmd/home

# Re-generate templates and SQL when relevant files change.
#
# Note that we can't include this as a part of "web" because if we do, when go
# generate runs, it changes the modification time of the files even if they
# haven't changed, which causes the task to be restarted, and go generate to run
# again. Our workaround is to only run go generate when files need to be
# regenerated, and avoid watching on go generate input files in web.
assets: reflex --decoration=none -r '\.(templ|sql)$' -- go generate ./...
