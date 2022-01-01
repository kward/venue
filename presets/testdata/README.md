# README for Presets Testdata

- -----------------------------------------------------------------------------
## D-Show 4 Band EQ [44455134]

211231.00 Clean EQ / clean EQ with all values zeroed, and the EQ is off.

- -----------------------------------------------------------------------------
## D-Show Input Channel (Built-In > VENUE Input Channel)

- Strings are \0 terminated.
- Strip = channel name

### 211231.00 Ch 1 Clear Console

Loaded the "Demo Shows > Clear Console" show.

- strip: Ch 1

### 211231.01 Ch 1 again

Saved the Ch 1 preset again with no changes in the interface.

- strip: Ch 1

```
$ cmp -lb 211231.00* 211231.01*
```

No changes :-)

### 211231.02 Ch 2

Saved a Ch 2 preset with no interface changes.

- strip: Ch 2

```
$ cmp -lb 211231.00* 211231.02*
3061  61 1     62 2
3080   0 ^@     1 ^A
```

### 211231.02 Ch 9

Saved a Ch 9 preset with no interface changes.

- strip: Ch 9

```
$ cmp -lb 211231.00* 211231.03*
3061  61 1     71 9
3080   0 ^@    10 ^H
```

The value at position 3080 seems to indicate the input position. Also, the `cmp` command is outputting in octal.

```
$ diff -u <(xxd 211231.00*) <(xxd 211231.03*)
--- /dev/fd/63	2021-12-31 16:21:00.000000000 +0100
+++ /dev/fd/62	2021-12-31 16:21:00.000000000 +0100
@@ -189,5 +189,5 @@
 00000bc0: ffff 000a 0000 0000 204e 0000 0000 0000  ........ N......
 00000bd0: 0010 0020 4e00 0064 0000 0064 0000 000a  ... N..d...d....
 00000be0: 5374 7269 7000 0d0b 0000 0000 0060 faff  Strip........`..
-00000bf0: ff43 6820 3100 0a53 7472 6970 2054 7970  .Ch 1..Strip Typ
-00000c00: 6500 0d02 0000 0000 01                   e........
+00000bf0: ff43 6820 3900 0a53 7472 6970 2054 7970  .Ch 9..Strip Typ
+00000c00: 6500 0d02 0000 0008 01                   e........
```

Looking closer at the start of the file, all strings appear to be prefixed with a `\0a` (newline) character. This could be useful.
