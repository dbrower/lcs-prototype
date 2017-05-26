lcs-prototype
-------------

Proof-of-concept of a way to manage licensed content in the library.


Data Model
----------

There are two types of items: _resources_ and _files_.
A _resource_ is akin to an item level record.
A _file_ is binary content that can be downloaded.
A resource can have any number of files, including zero.
Conversely, every file is part of exactly one resource.

All information is collected by scanning a file tree.
Each resource is stored in its own directory.
A resource is indicated when a directory named _X_ contains a file named _X-meta.xml_.
Every file inside that directory and inside any contained directories is then part of that resource.
One resource cannot be stored inside another resource's directory.

Resource directories can be organized in any other way.
For example we could have a directory tree like the following

    er/
      eh_online/
        eh_online-meta.xml
        a.pdf
        map1/
          map1.bst
          map1.gif
      other/
        blah/
          blah-meta.xml
      lapop/
        lapop_2017/
          lapop_2017-meta.xml
          file1.csv
          file2.csv

This file tree contains three resources, `eh_online`, `blah`, and `lapop_2017`.
The `eh_online` resource contains four files.

    eh_online-meta.xml
    a.pdf
    map1/map1.bst
    map1/map1.gif

The `blah` resource contains one file: `blah-meta.xml`, and the `lapop_2017` resource contains three files.

The meta.xml file is an XML file containing the metadata for the item.
It is a top level tag of `<metadata>` containing a bunch of arbitrary tags---essentially a bunch of key-value pairs (that can be repeated).
Any key can be present. There are a few predefined keys, though

 * `identifier` - the id for this resource. Must be the same as the directory name.
 * `title` - the title to use for the web page.
 * `addeddate` - the date this item was added to the system (iso8660).
 * `accesslevel` - the permissions needed to access this item and all the files it contains.
 * `status` - the status of this resource...one of "active", "offline", or "tombstone". Defaults to "active".
 * `redirect` - if present, a URL to do a 301 (permanent) redirect to.

The `accesslevel` levels is one of the following.

 * `public` - anyone may view this
 * `login` - anyone who is logged in to the site may view this
 * `private` - no one may view this

Every file in a resource has the same access level as the resource itself.

The status levels mean the following:

 * `active` - this resource is available
 * `offline` - this resource is temporarily unavailable
 * `tombstone` - this resource is deleted

An example metadata file is:

    <metadata>
      <identifier>lapop_2017</identifier>
      <alephnumber>0001234567</alephnumber>
      <title>Latin America Population Data, 2017</title>
      <accesslevel>login</accesslevel>
      <addeddate>2017-04-21T15:06:45Z</addeddate>
    <metadata>

The system will create and maintain a FILES.xml file containing a list of all the files and their checksums.
The system will also create and maintain a LOG.txt file containing a provenance log for the resource.

The user interacts with the system using a Web interface.
Each resource is identified with the route `/{resource_id}`.
Each file can be downloaded with the route `/{resource_id}/{file_name}`.

The following logic is followed when serving requests:

 1. Does this resource exist? If not, return 404.
 1. Does the user have permission to view this resource? If no, say so and (possibly) prompt to log in.
 1. Is there a redirect URL? If so, redirect
 1. Is the status offline or tombstone? If so, display a page stating the resource is unavailable or that it is gone.
 1. If a file is given, does it exist? If not, return 404.
 1. If a file is given, return the file contents
 1. If there is a `page.html` file, display that
 1. Display a directory listing for the resource, with option to download a zip file.


