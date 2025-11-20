# TipTap Editor Toolbar Changes

## Date: November 20, 2025

## Summary
Added a visual toolbar to the TipTap editor with formatting buttons.

## Files Modified
- `/home/conor/dev/wiki-notes/frontend/src/lib/components/Editor.svelte`

## Changes Made

### 1. Added State Variables (Line ~24-36)
```typescript
// Changed from:
let editorElement: HTMLDivElement;
let editor: Editor | null = null;

// To:
let editorElement = $state<HTMLDivElement>();
let editor = $state<Editor | null>(null);

// Added new state for tracking active formatting:
let isActive = $state({
    bold: false,
    italic: false,
    strike: false,
    code: false,
    heading: { 1: false, 2: false, 3: false },
    bulletList: false,
    orderedList: false,
    codeBlock: false,
    blockquote: false
});
```

### 2. Added Editor Event Handlers (Line ~72-77)
```typescript
// Modified onUpdate and added onSelectionUpdate:
onUpdate: ({ editor }) => {
    content = editor.getHTML();
    updateActiveStates();  // NEW
},
onSelectionUpdate: () => {   // NEW
    updateActiveStates();    // NEW
}
```

### 3. Added updateActiveStates Function (Line ~86-104)
```typescript
function updateActiveStates() {
    if (!editor) return;
    isActive = {
        bold: editor.isActive('bold'),
        italic: editor.isActive('italic'),
        strike: editor.isActive('strike'),
        code: editor.isActive('code'),
        heading: {
            1: editor.isActive('heading', { level: 1 }),
            2: editor.isActive('heading', { level: 2 }),
            3: editor.isActive('heading', { level: 3 })
        },
        bulletList: editor.isActive('bulletList'),
        orderedList: editor.isActive('orderedList'),
        codeBlock: editor.isActive('codeBlock'),
        blockquote: editor.isActive('blockquote')
    };
}
```

### 4. Added Toolbar HTML (Line ~195-315)
Added a complete toolbar div with buttons for:
- Bold, Italic, Strikethrough, Inline Code
- H1, H2, H3 headings
- Bullet List, Numbered List
- Code Block, Blockquote
- Undo, Redo

Each button has SVG icons and active state styling.

### 5. Added Toolbar CSS Styles (Line ~391-426)
```css
.toolbar-btn {
    /* Button styling with transitions */
}

.toolbar-btn:hover:not(:disabled) {
    /* Hover effects */
}

.toolbar-btn.active {
    /* Blue highlight for active formatting */
}
```

## Potential Performance Issues

The lag could be caused by:
1. **Frequent state updates** - `updateActiveStates()` is called on every selection change and content update
2. **Multiple `isActive()` checks** - We check 11+ different formatting states on every cursor movement
3. **Re-rendering toolbar** - The entire toolbar re-renders when `isActive` state changes
4. **SVG rendering** - 15+ SVG icons being re-evaluated on each render

## To Revert
Simply restore the previous version of `Editor.svelte` from git:
```bash
git checkout HEAD -- frontend/src/lib/components/Editor.svelte
```

Or manually remove:
- The `isActive` state variable
- The `updateActiveStates()` function
- The `onSelectionUpdate` handler
- The toolbar div (everything between `<!-- Toolbar -->` and `<!-- Editor -->`)
- The `.toolbar-btn` CSS rules
- Change `let editor` and `let editorElement` back to non-reactive declarations

