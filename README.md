# fexcel

Utility for setting FANUC robot comments from an Excel spreadsheet.

## Usage

    fexcel [options] filename host

    > ./fexcel -sheet Data -numregs A2 -posregs D2 spreadsheet.xlsx 127.0.0.101

## Options

| Option   | Description |
| -------- | ----------- |
| -sheet   | the name of the sheet (default Sheet1) |
| -offset  | number of columns to shift between id and comments (default 1) |
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

## Details

fexcel assumes that your spreadsheet has indices for a given item that start
in the provided cell, and the comments for that item are in the column you
provided plus the offset flag (default 1).

e.g. in the above example, the numeric register ids start in cell A2 with
comments starting in cell B2. Position registers ids start in cell D2 with
comments starting in E2.

If you have data spread out on different sheets, you will have to run fexcel
multiple times: once for each sheet.
