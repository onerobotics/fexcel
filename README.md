# fexcel

Utility for setting FANUC robot comments from an Excel spreadsheet.

Download the latest release [here](https://github.com/onerobotics/fexcel/releases/latest).

## Usage

Make sure KAREL is unlocked under `Setup > Host Comm > HTTP`.

    fexcel [options] filename host(s)...

    > ./fexcel -sheet Data -numregs A2 -posregs D2 -dins IO:A2 spreadsheet.xlsx 127.0.0.101 127.0.0.102

## Options

| Option   | Description |
| -------- | ----------- |
| -sheet   | the name of the default sheet (default Sheet1) |
| -offset  | number of columns to shift between id and comments (default 1) |
| -timeout | how long to wait for robot(s) to respond in milliseconds (default 500) |
| -noupdate| skip the check for an updated version of fexcel |
| -numregs | start cell of numeric register definitions |
| -posregs | start cell of position register definitions |
| -ualms   | start cell of user alarm definitions | 
| -rins    | start cell of robot input definitions |
| -routs   | start cell of robot output definitions |
| -dins    | start cell of digital input definitions |
| -douts   | start cell of digital output definitions |
| -gins    | start cell of group input definitions |
| -gouts   | start cell of group output definitions |
| -ains    | start cell of analog input definitions |
| -aouts   | start cell of analog output definitions |
| -sregs   | start cell of string register definitions |
| -flags   | start cell of flag definitions |

Note that start cell flags can be optionally prefixed with a sheet name that
overrides the default `-sheet` parameter. (e.g. `-numregs Data:A2`)

## Details

fexcel assumes that your spreadsheet has indices for a given item that start
in the provided cell, and the comments for that item are in the column you
provided plus the offset flag (default 1).

e.g. in the above example, the numeric register ids start in cell A2 with
comments starting in cell B2. Position registers ids start in cell D2 with
comments starting in E2. Digital input ids start in cell A2 on the IO sheet.
