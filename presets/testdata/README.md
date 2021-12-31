# README for Presets Testdata

## D-Show 4 Band EQ [44455134]

211231.00 Clean EQ / clean EQ with all values zeroed, and the EQ is off.

## D-Show Input Channel (Built-In > VENUE Input Channel)

- strip = channel name

### 211231.00 Reset Ch 57
channel reset
- console: S6L 24C
- strip: Ch 57
- channel color: Green
- input position: 65

### 211231.01 Reset Ch 58
Channel reset with right-click > Reset Channel
- console: S6L 24C
- strip: Ch 58
- channel color: Green
- inputs position: 66

```
$ cmp -lb 211231.0[01]*
3100  67 7     70 8
3119  70 8     71 9
```

Byte #3100 changed from '7' to '8' as expected.
Byte #3119 changed unexpectedly. 

```
$ xxd 211231.00\ Reset\ Ch\ 57.ich |tail -3
00000c00: 0064 0000 000a 5374 7269 7000 0d0c 0000  .d....Strip.....
00000c10: 0000 0000 0000 0043 6820 3537 000a 5374  .......Ch 57..St
00000c20: 7269 7020 5479 7065 000d 0200 0000 3801  rip Type......8.
```

Does byte #3119 indicate a different "strip type" ?

### 211231.02 Reset Ch 58 2
Preset saved again with no changes.
- console: S6L 24C
- strip: Ch 58
- channel color: Green
- inputs position: 66

```
$ cmp -lb 211231.0[12]*
```

There was no diff, which indicates that the time is not stored in the file.

### 211231.03 Reset Ch 57 2
Ch 58 channel name manually changed to "Ch 57".
- console: S6L 24C
- strip: Ch 58
- channel color: Green
- inputs position: 66

```
$ cmp -lb 211231.0[23]*
3100  70 8     67 7
```

As expected, only byte #3100 changed.
