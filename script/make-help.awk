#!/usr/bin/awk -f
#
# Prints target help comments (##) included in Makefiles.
# Loosely based on: https://gist.github.com/prwhite/8168133#gistcomment-2278355

/^[\-\_[:alnum:]]+:/ {
    comment = match(previous_line, /^## (.*)/);
    if (comment) {
        target = substr($1, 0, index($1, ":") - 1);
        comment = substr(previous_line, RSTART + 3, RLENGTH);
        printf "  %-21s %s\n", target, comment;
    }
}
{ previous_line = $0 }
