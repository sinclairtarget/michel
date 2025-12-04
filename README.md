# Michel
![brutalist](https://github.com/user-attachments/assets/80490b07-8eb6-4a6a-82d7-185c0964a1df)

Michel is a minimal static site builder based on [MyST Markdown](https://mystmd.org/).

Michel favors guts-visible explicitness to "convention over configuration."

## Model
Michel builds a site by reading input files from these directories:

`content`
: Your website content / prose, written using MyST Markdown.

`site`
: Your website HTML pages (templated using Go templating) and assets.

`layouts`
: Your templated layouts that can be shared among multiple pages.

`partials`
: Your templated components that can be shared among multiple pages.

After processing, all output gets written to the target directory named
`public`.

Michel never automatically creates any pages for you. Every page in your
website must exist as a page under `/site`. No content defined under
`/content` appears in your website unless it is explicitly used by a page
under `/site`.
