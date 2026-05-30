/*
Copyright (c) 2023-2026 Microbus LLC and various contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mermaid

// Color constants resolve to BESPA Material Design 3 palette CSS variable
// references, suitable for passing as any color value to a Mermaid diagram
// (classDef, style, linkStyle). The widget's var-expansion bridges them
// through Mermaid's parser at render time; the diagram tracks the host page's
// light/dark theme without re-rendering.
//
// Use these constants instead of writing rgb(var(--md-sys-color-X)) by hand:
// the underlying CSS token names are an implementation detail of BESPA and
// may change. A caller that wants a token not in this set can still write a
// raw var() reference — the var-expansion in varexpand.go handles arbitrary
// CSS custom properties.
//
// Names mirror the Material Design 3 system color roles. Each "X" surface
// color pairs with an "OnX" text color chosen for contrast.
const (
	Primary            = "rgb(var(--md-sys-color-primary))"
	OnPrimary          = "rgb(var(--md-sys-color-on-primary))"
	PrimaryContainer   = "rgb(var(--md-sys-color-primary-container))"
	OnPrimaryContainer = "rgb(var(--md-sys-color-on-primary-container))"

	Secondary            = "rgb(var(--md-sys-color-secondary))"
	OnSecondary          = "rgb(var(--md-sys-color-on-secondary))"
	SecondaryContainer   = "rgb(var(--md-sys-color-secondary-container))"
	OnSecondaryContainer = "rgb(var(--md-sys-color-on-secondary-container))"

	Tertiary            = "rgb(var(--md-sys-color-tertiary))"
	OnTertiary          = "rgb(var(--md-sys-color-on-tertiary))"
	TertiaryContainer   = "rgb(var(--md-sys-color-tertiary-container))"
	OnTertiaryContainer = "rgb(var(--md-sys-color-on-tertiary-container))"

	Error            = "rgb(var(--md-sys-color-error))"
	OnError          = "rgb(var(--md-sys-color-on-error))"
	ErrorContainer   = "rgb(var(--md-sys-color-error-container))"
	OnErrorContainer = "rgb(var(--md-sys-color-on-error-container))"

	// Ok is BESPA's success extension to the MD3 palette, used for affirmative
	// states (validated input, healthy status, etc.). Not part of the MD3 spec.
	Ok            = "rgb(var(--md-sys-color-ok))"
	OnOk          = "rgb(var(--md-sys-color-on-ok))"
	OkContainer   = "rgb(var(--md-sys-color-ok-container))"
	OnOkContainer = "rgb(var(--md-sys-color-on-ok-container))"

	Background   = "rgb(var(--md-sys-color-background))"
	OnBackground = "rgb(var(--md-sys-color-on-background))"

	Surface          = "rgb(var(--md-sys-color-surface))"
	OnSurface        = "rgb(var(--md-sys-color-on-surface))"
	SurfaceVariant   = "rgb(var(--md-sys-color-surface-variant))"
	OnSurfaceVariant = "rgb(var(--md-sys-color-on-surface-variant))"
	SurfaceDim       = "rgb(var(--md-sys-color-surface-dim))"
	SurfaceBright    = "rgb(var(--md-sys-color-surface-bright))"
	SurfaceTint      = "rgb(var(--md-sys-color-surface-tint))"

	SurfaceContainerLowest  = "rgb(var(--md-sys-color-surface-container-lowest))"
	SurfaceContainerLow     = "rgb(var(--md-sys-color-surface-container-low))"
	SurfaceContainer        = "rgb(var(--md-sys-color-surface-container))"
	SurfaceContainerHigh    = "rgb(var(--md-sys-color-surface-container-high))"
	SurfaceContainerHighest = "rgb(var(--md-sys-color-surface-container-highest))"

	Outline        = "rgb(var(--md-sys-color-outline))"
	OutlineVariant = "rgb(var(--md-sys-color-outline-variant))"

	InversePrimary   = "rgb(var(--md-sys-color-inverse-primary))"
	InverseSurface   = "rgb(var(--md-sys-color-inverse-surface))"
	InverseOnSurface = "rgb(var(--md-sys-color-inverse-on-surface))"

	Scrim  = "rgb(var(--md-sys-color-scrim))"
	Shadow = "rgb(var(--md-sys-color-shadow))"
)
