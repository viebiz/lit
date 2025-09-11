# Localization

This project uses [go-i18n](https://github.com/nicksnyder/go-i18n) to translate user-facing messages.

## Message bundles

Create a message bundle with `i18n.Init` and point it to the directory containing your translation files. Bundles lazily load a file the first time a language is requested:

```go
bundle := i18n.Init(ctx, i18n.BundleConfig{
    DefaultLang:      "en",
    SourcePath:       "resources/i18n",
    BundleFileFormat: "json",
})

// Vietnamese messages are loaded on demand
lc := bundle.GetLocalize("vi")
msg := lc.Localize("hello", nil)
```

`LoadMessageFile` can be called manually to preload a language or to support formats beyond JSON.

## i18n helpers

- `i18n.Localizable` exposes `Localize` and `TryLocalize` for translating message IDs.
- `i18n.SetInContext` and `i18n.FromContext` store a localizer in `context.Context` so any part of a request can access it.

## Context aware localizer

`middleware/http.LocalizationMiddleware` selects the language from the `Accept-Language` header (falling back to a default) and injects a localizer into the request context. Handlers can then obtain it with `i18n.FromContext`.

```go
r.Use(http.LocalizationMiddleware(ctx, http.Config{}))
```

## Translating validation errors

```go
func createUser(c lit.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        lc := i18n.FromContext(c.Request().Context())
        msg := lc.Localize("validation.invalid_payload", nil)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
    }

    if err := validate.Struct(req); err != nil {
        lc := i18n.FromContext(c.Request().Context())
        for _, verr := range err.(validator.ValidationErrors) {
            fieldMsg := lc.Localize(verr.Tag(), map[string]interface{}{"field": verr.Field()})
            // handle translated fieldMsg ...
        }
    }

    // ...
    return nil
}
```

## Supported languages

English is supported out of the box. Add another language by creating a `<lang>.json` file in `resources/i18n` (e.g. `vi.json`) and placing the translated messages inside. The bundle will discover it automatically the first time the language key is requested.

