# tui-pipes

Create dashboards that communicate by text.

```
 ┌────────────────────────┐
 │                        │       ┌─────────────────────────────────────────────────────┐
 │      navbar            │       │  name                                               │
 │                        │       │                                                     │
 └────────────────────────┘       │  cmd:                                               │
 ┌───────────┐┌───────────┐       │                                                     │
 │           ││           │       │  args:                                              │
 │   list    ││   preview │       └─────────────────────────────────────────────────────┘
 │           ││           │
 │           ││           │
 │           ││           │
 │           ││           │
 │           ││           │
 │           ││           │
 │           ││           │
 └───────────┘└───────────┘

```

Inspired by lazygit, fzf.

Necessary background:

1. `tview`
   - focus handling
     - primitive
     - box
     - grid
     - form
   - draw
     - box
     - grid
     - form
2. ANSI writer
3. `context` package

## Bad decisions

Navbar as separate component
Using tview: just use tcell instead, use our own event loop
App as separate component: it is just clutter for now, but maybe later pages may be useful

## Usage examples

1. Show files

- filename - `show_files.json`
- source - `find [., -type, f]`
- preview - `bat $item`
- bindings
  - `Enter`
  - `Space`
  - `?`
