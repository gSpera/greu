# Go Roff Extension Util
greu is an utility to quickly enhance troff with custom commands

## Usage
greu is yet another preprocessor to roff, you can use greu in any position in your pipeline but you want to use as the first one, for example:

```sh
cat input.roff | greu | eqn | roff > output.ps
```

## Commands
greu defines the commands in a config, the primary fields are:

| Name | Meaning | Example |
-------|---------|-------|
| Cmd  | The shell command to execute | gnuplot | 
| OpenTag | The opening tag detected in the input file by greu | .GNUPLOTSTART |
| CloseTag | The closing tag detected in the input file by greu | .GNUPLOTEND | 
| ReplaceOpenTag | A string that will be inserted in the output before the command output | .PS |
| ReplaceCloseTag | A string that will be inserted in the output after the command output | .PE |
| InputPrefix | A string that will be passed to the command before the input file lines, usefull for some prelude | set terminal eps | 
| InputPostfix | A string that will be passed to the command before closing, usefull to exit | exit | 

## Example
Let's say that we want to incorportate gnuplot graphs in our document and that our roff implementation has a .PS/.PE commands to include raw eps.

```
Cmd: gnuplot
OpenTag: .GNUPLOTSTART
CloseTag: .GNUPLOTEND
ReplaceOpenTag: .PS
ReplaceCloseTag: .PE
InputPrefix: set terminal eps
InputPostfix: exit
```

Now we can write our document:
```
Commodo excepteur fugiat non deserunt laboris nisi culpa sit.
.GNUPLOTSTART
plot sin(x)
.GNUPLOTEND
Do consectetur laborum minim exercitation amet minim sit laboris adipisicing exercitation proident labore proident labore.
```
After passing in greu the output will be:
```
Commodo excepteur fugiat non deserunt laboris nisi culpa sit.
.PS
%!PS-Adobe-3.0 EPSF-3.0
%%Creator: cairo 1.17.6 (https://cairographics.org)
%%CreationDate: Sun Dec 25 13:17:18 2022
%%Pages: 1
%%DocumentData: Clean7Bit
%%LanguageLevel: 2
%%BoundingBox: 0 0 360 216
%%EndComments
%%BeginProlog
...
.PE
Do consectetur laborum minim exercitation amet minim sit laboris adipisicing exercitation proident labore proident labore.
```

## How it works
When the open tag is found greu executes the command, greu continues reading and the file is passed in the stdin of the command, the stdout is then written to the output