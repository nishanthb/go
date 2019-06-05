package main

import(
"fmt"
"regex"
)

func main() { 
var tomatch string = "fwgfwfdsfkwfuwfewuefewfoo=ksweuf48ffsd=kfwefgateway=ekffefwfranger=effwefnwf"
 r := regexp.MustCompile(`foo=(?P<foo>[^ ]+)\s+gateway=(?P<gw>[^ ]+)\s+ranger=(?P<rngr>[^ ]+)\s`)
    if r.MatchString(tomatch) == true {
        results := findNamedMatches(r, tomatch)
        fmt.Printf("%#v\n", results)
    }
}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
    match := regex.FindStringSubmatch(str)

    results := map[string]string{}
    for i, name := range match {
        results[regex.SubexpNames()[i]] = name
    }
    return results
}
