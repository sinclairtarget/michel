# Michel
![brutalist](https://github.com/user-attachments/assets/80490b07-8eb6-4a6a-82d7-185c0964a1df)

Michel is a minimal static site builder based on [MyST Markdown](https://mystmd.org/).

## Why MyST?
MyST is a Markdown specification that adds several useful features to Markdown
in a standardized way, including tables, footnotes, inline math, and
[admonitions](https://mystmd.org/spec/admonitions).

MyST also standardizes a syntax for extending the format with custom
"directives" and "roles." Michel provides a small set of custom directives and
roles for use in Markdown content files. These provide functionality similar to
Hugo's shortcodes without requiring any templating of content files.

Finally, if you want to publish your writing in other ways, you aren't
restricted to HTML. Tools in the wider MyST ecosystem can turn your content
files into Word documents, PDFs, and more. 

## Build Overview
Michel builds a site by reading input files from these directories:

`content`: Your website content / prose, written using MyST Markdown.

`site`: Your website HTML pages (templated using Go templating) and assets.

`layouts`: Your templated layouts that can be shared among multiple pages.

`partials`: Your templated sub-components that can be shared among multiple pages.

After processing, all output gets written to the target directory named
`public`.

### The "Brutalist" Part
Michel never infers the existence of any page based on your content files. The
organization of your content under the `content` directory does not imply
anything about the organization of pages in your built site. No content will
appear in your website unless it is explicitly used by a page template under
the `site` directory. No page will appear in your website unless a page
template exists for that path under `site`.
