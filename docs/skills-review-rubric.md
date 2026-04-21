# Skills Review Rubric — zesh-one-skills

This rubric is the review standard for every skill in `skills/**/SKILL.md`.

It evaluates each skill across three layers:

1. **Spec Compliance** — mandatory
2. **Agent Skills Best Practices** — quality and context efficiency
3. **Repo Philosophy Fit** — whether the skill captures the project's specific "seasoning" instead of generic fundamentals

---

## 1) Layer 1 — Spec Compliance

If a skill fails this layer, it is **blocked** until corrected.

### 1.1 Frontmatter

#### `name`
- [ ] Present
- [ ] 1-64 characters
- [ ] Lowercase letters, numbers, and hyphens only
- [ ] Does not start or end with `-`
- [ ] Does not contain consecutive hyphens (`--`)
- [ ] Matches the parent directory name exactly

#### `description`
- [ ] Present
- [ ] Non-empty
- [ ] Describes what the skill does
- [ ] Describes when to use it
- [ ] Contains useful trigger keywords

#### `license`
- [ ] If present, valid and concise

#### `compatibility`
- [ ] If present, adds real environment requirements
- [ ] Not used as a miscellaneous note field

#### `metadata`
- [ ] If present, uses coherent keys
- [ ] Does not duplicate spec-defined fields unnecessarily

#### `allowed-tools`
- [ ] If present, uses a single space-separated string
- [ ] Not written as a YAML list

### 1.2 Minimum structure
- [ ] `SKILL.md` exists
- [ ] File is valid YAML frontmatter followed by Markdown

### Layer 1 result
- **PASS**
- **FAIL** → fix spec issues before content cleanup

---

## 2) Layer 2 — Agent Skills Best Practices

This layer checks whether the skill is well-calibrated for real agent use.

### 2.1 Scope and activation
- [ ] Covers a coherent unit of work
- [ ] Not too broad
- [ ] Not so fragmented that many skills must load for one task
- [ ] Description activates the skill clearly without excessive ambiguity

#### Signals of trouble
- Mixes multiple domains or responsibilities
- Combines implementation, operations, architecture, audit history, and reference material in one file
- Should likely be split into `SKILL.md` + `references/`

### 2.2 Context economy
- [ ] Adds knowledge the agent would not know by default
- [ ] Avoids generic fundamentals
- [ ] Avoids empty "best practices" language
- [ ] Avoids duplicating framework or library documentation

#### Core test
Ask for each section:

> Would the agent likely get this wrong without this instruction?

If the answer is **no**, the section is a candidate for removal.

### 2.3 Defaults and calibration
- [ ] Provides a clear default approach
- [ ] Mentions alternatives only as escape hatches
- [ ] Avoids presenting long menus of equally weighted options
- [ ] Is prescriptive where the task is fragile
- [ ] Leaves flexibility where multiple valid approaches exist

### 2.4 Progressive disclosure
- [ ] `SKILL.md` contains the operational core
- [ ] Detailed material is moved to `references/`, `assets/`, `scripts/`, or support files
- [ ] References tell the agent when to read each extra file
- [ ] Historical notes do not dominate the main operational file

#### Recommended thresholds
- Ideal: **80-140 lines**
- Acceptable: **up to ~180-220 lines**
- Risk: **over 220 lines**
- Strong smell: **over 500 lines**

These thresholds do **not** automatically mean spec failure. They are design-quality indicators.

### 2.5 Reusability of instructions
- [ ] Teaches a reusable procedure, not just a one-off answer
- [ ] Includes real gotchas where valuable
- [ ] Includes constraints or tradeoffs when relevant
- [ ] Uses examples sparingly and only where they add clarity

### Layer 2 result
- **GOOD**
- **TRIM**
- **REWRITE**

---

## 3) Layer 3 — Repo Philosophy Fit

This layer captures the editorial standard of this repository.

### 3.1 Is it actually a skill?
- [ ] Encodes team-specific conventions
- [ ] Encodes real constraints
- [ ] Encodes anti-patterns seen in practice
- [ ] Encodes exceptions and tradeoffs
- [ ] Does not read like a generic handbook

### 3.2 The chef test
Ask:

> Is this telling the chef how to chop vegetables, or how *this kitchen* wants the dish prepared?

If it mostly teaches fundamentals, it fails the philosophy fit.

### 3.3 Signals of misalignment
- [ ] Too much general theory
- [ ] Too much framework onboarding
- [ ] Too many long snippets
- [ ] Embedded changelog or history creates noise
- [ ] Mixes runbook, tutorial, handbook, and skill in one place
- [ ] Feels encyclopedic instead of operational

### Layer 3 result
- **ALIGNED**
- **DRIFTING**
- **MISALIGNED**

---

## 4) Final review template

Use this structure when reviewing any skill:

```md
## <skill-name>

### Layer 1 — Spec Compliance
- Result: PASS / FAIL
- Issues:
  - ...

### Layer 2 — Best Practices
- Result: GOOD / TRIM / REWRITE
- Issues:
  - ...

### Layer 3 — Repo Philosophy Fit
- Result: ALIGNED / DRIFTING / MISALIGNED
- Issues:
  - ...

### Final Decision
- KEEP
- TRIM
- REWRITE
- BLOCKED (fix spec first)
```

---

## 5) Decision rules

### BLOCKED
Use when the skill fails spec compliance.

### KEEP
Use when the skill:
- Passes spec
- Is good according to best practices
- Is aligned with repo philosophy

### TRIM
Use when the skill passes spec but has excess noise, duplication, overly long examples, or unnecessary history.

### REWRITE
Use when the skill may partially pass spec but is conceptually mis-scoped, poorly calibrated, or functions more like a handbook than a skill.

---

## 6) Fast operational checklist

- [ ] Does it pass spec?
- [ ] Does `name` match the directory?
- [ ] Is `allowed-tools` formatted correctly?
- [ ] Does it explain things the agent already knows?
- [ ] Does it give clear defaults?
- [ ] Does it include real gotchas?
- [ ] Is `SKILL.md` carrying too much context?
- [ ] Should some material move to `references/`?
- [ ] Does it feel like a skill or a manual?

---

## 7) Governing principle

> A skill must contribute specific, reusable, hard-to-infer knowledge using the minimum context necessary.

This principle combines:
- Agent Skills spec compliance
- Agent Skills best-practice quality
- The editorial philosophy of `zesh-one-skills`
