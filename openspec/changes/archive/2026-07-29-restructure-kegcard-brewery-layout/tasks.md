## 1. Markup restructure

- [x] 1.1 Add a new `brewery-plaque` block to `KegCard.svelte`, conditionally rendered when `keg.brewery` is non-empty, containing the brewery logo (`hasBreweryImage`) and brewery name
- [x] 1.2 Remove the existing `.brewery-row` (inline logo + name) from inside `.info`
- [x] 1.3 Reorder `.inner`'s children to `[brewery-plaque?] [info] [stats]`, keeping the existing `.divider` between info and stats

## 2. Styling

- [x] 2.1 Style `.brewery-plaque` as a fixed-width, `flex-shrink: 0` column (110-140px logo + name, centered), with a background/border treatment distinct from the existing hairline `.divider`
- [x] 2.2 Size and center the brewery logo image within the plaque; verify no layout shift when the logo is absent (name-only plaque)
- [x] 2.3 Add wrap/truncation handling for long brewery names within the fixed plaque width
- [x] 2.4 Verify `.info` and `.stats` reflow correctly when the plaque is omitted (no brewery) — no leftover gap or misalignment

## 3. Verification

- [x] 3.1 Manually verify on the carousel with: a keg with brewery logo + name, a keg with brewery name only (no logo), and a keg with no brewery at all
- [x] 3.2 Verify beer image size/shape variants (circle/square/can) are unaffected by the restructure
- [x] 3.3 Check layout at representative TV width and a smaller laptop width (e.g. 1366px) for crowding
- [x] 3.4 Run `npm run check` and `npm run lint` in `web/`
