# Scope, architecture and language menus
Branch: feat/zesh-skills-wizard-languages; baseline 32ee027200b56cff12ac7fca3786ef8673725df6. Delivery: local work-unit commits only; no source push, PR, release or global installation requested.
Forecast: ~900 authored diff lines. Checks: focused/full Go tests, vet, Node suite, six CGO-disabled builds; independent verification follows native assessment.

## Specs
- S1: "ok quedo bien, agrega la misma seleccion para arquitectura y para scope"
- S2: "adicional quiero que expandas el catalogo de arquitecturas con mas tipos y tipos hibridos comunes, con su breve explicacion"
- S3: "quiero una opcion al inicio de la instalacion o un parametro para elegir entre español o ingles (default ingles) talvez con un parametro basta algo asi\nzesh-skills --es"
- S4: "la CLI debe ser multiplataforma, un binario sin dependencias"
- S5: "las skills solo se instalan localmente en el proyecto a trabajar"

## Tasks
- T1 | S1-S5 | inline exploration/preparation | done | Previous TUI manually accepted; baseline local commit 32ee027; bounded error/protocol mapping read back.
- T2 | S1-S2 | inline (background writers disabled) | in_progress | Scope multi-select, architecture single/custom, expanded explained profiles, tests/docs; commit pending.
- T3 | S3-S5 | inline (same parser/picker unit) | pending | --es presentation boundary, default English compatibility, localized scope/profile UI and diagnostics; commit pending.
- T4 | S1-S5 | native-assessment plan + parent checks | pending | Observe checks/builds and required independent proof; never reopen old winpty probe; commit pending.

## Log
L1 User request (verbatim):
> ok quedo bien, agrega la misma seleccion para arquitectura y para scope
> adicional quiero que expandas el catalogo de arquitecturas con mas tipos y tipos hibridos comunes, con su breve explicacion
> quiero una opcion al inicio de la instalacion o un parametro para elegir entre español o ingles (default ingles) talvez con un parametro basta algo asi
> zesh-skills --es
L2 Previous TUI manual validation closes the prior terminal-exit blocker. Previous work captured locally as 32ee027; no CLI source was published. Catalog main remains 9b099b0. Existing terminal raw/restore adapters and pinned dependencies will not be redesigned.
L3 Technical choices under existing autonomy: reuse a synchronous generic picker; scope accepts several domains, architecture chooses one profile and retains Custom text. Central local architecture catalog includes simple/layered/vertical-slice/Clean/Hexagonal/Onion/modular/distributed/UI styles and common hybrids; descriptions distinguish domain organization, deployment styles and patterns. Selection stores a stable profile ID, not a translated label, and creates no infrastructure or mandatory skill bundle.
L4 --es is the sufficient requested language mechanism; no extra startup question or persisted locale. English remains byte-compatible by default. Localize owned prompts, menus, help, preview, confirmations and diagnostics at presentation boundaries. Do not translate IDs, flags, JSON keys, hashes, custom user values, downloaded skill documents, original catalog descriptions or generated agent instruction formats. OS/technical details may retain their original language, explicitly identified in Spanish.
L5 Explorer evidence: English errors also originate in catalog.go/install.go, including a modified-index confirmation suffix. Existing verification_test.go asserts exact default-English preview/error streams and a three-read conflict fixture. Preserve these by leaving domain/storage code unchanged and localizing UI templates before inserting opaque values.
L6 Test-first: each new menu/language rule gets a RED public-interface/model test; changed shared picker must retain old skill key behavior and automation flags. Scope/architecture cancellation must not fetch or write; --help/--es --help need no network. Native assessment decides outside verification/review, never a TODO transition.
