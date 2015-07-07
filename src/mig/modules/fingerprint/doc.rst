========================================
Mozilla InvestiGator: Fingerprint module
========================================
:Author: Aaron Meihm <ameihm@mozilla.com>

.. sectnum::
.. contents:: Table of Contents

The fingerprint module identifies and returns information that matches a
set of known fingerprints specified in the module itself. This can be used
to for example, return version information associated with arbitrary installed
libraries, running kernel information, loaded kernel modules, and so on.

Usage
-----

From an agent perspective, the following actions are executed by the agent
using a given template when an action is received:

* The file system is scanned for any files that match the file name specification in the template
* For each of the matched files, they are scanned line by line for any content the matches the content expression
* For any matches, the subgroup specified in the content expression is obtained
* The subgroup is passed through the template transform function
* The results (file name, and final content string) is returned to the investigator

Included templates
~~~~~~~~~~~~~~~~~~

For a list of included templates, see the help output for the fingerprint
module.

.. code::

        $ bin/linux/amd64/mig fingerprint
        Query parameters
        ----------------
        -template <name>   - Scan using template
                           ex: template mediawiki
                           query for specific module supplied template

        -depth <int>       - Specify maximum directory search depth
                           ex: depth 2
                           default depth is 10

        -root <path>       - Specify search root
                           ex: root /usr/local
                           default root is /

        Available templates:

        linuxkernel         - Running Linux kernel information
        linuxmodules        - Loaded Linux modules
        pythonegg           - Python package versions
        django              - Django framework versions
        mediawiki           - MediaWiki framework versions

New templates can be added to the fingerprint module by updating `templates.go`
under `mig/modules/fingerprint/`.

Executing templates
~~~~~~~~~~~~~~~~~~~

To execute templates, the template to run is passed as an argument to
`template` with the fingerprint module. An example using the linuxmodules
template:

.. code::

        $ bin/linux/amd64/mig fingerprint -template linuxmodules
        1 agents will be targeted. ctrl+c to cancel. launching in 5 4 3 2 1 GO
        Following action ID 1436206513445352448.status=inflight.
        - 100.0% done in 3.40256105s
        1 sent, 1 done, 1 succeeded
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=ppdev
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=snd_intel8x0
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=snd_ac97_codec
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=serio_raw
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=ac97_bus
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=joydev
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=snd_pcm
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=i2c_piix4
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=snd_timer
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=snd
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=soundcore
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=parport_pc
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=parport
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=video
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=mac_hid
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=hid_generic
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=usbhid
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=hid
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=psmouse
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=ahci
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=libahci
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=e1000
        ubuntu-dev fingerprint name=linuxmodules root=/proc/modules entry=pata_acpi
        1 agent has found results

An example using the mediawiki template:

.. code::

        alm@ubuntu-dev:~/mig$ bin/linux/amd64/mig fingerprint -template mediawiki
        1 agents will be targeted. ctrl+c to cancel. launching in 5 4 3 2 1 GO
        Following action ID 1436207350958595328.status=inflight..status=completed
        - 100.0% done in 5.353326918s
        1 sent, 1 done, 1 succeeded
        ubuntu-dev fingerprint name=mediawiki root=/tmp/mediawiki/DefaultSettings.php entry=1.3.6
        1 agent has found results

Template structure
~~~~~~~~~~~~~~~~~~

Templates include the following attributes:

* **Filename**: The file name to locate on the file system. File names can be either the string to match against, or a regular expression. (**required**)
* **Content match**: A regular expression indicating the lines to match from the file. This regular expression must contain one subgroup; that subgroup is extracted as information to return to the investigator. (**required**)
* **Transform function**: A function to call to pass the extracted information from the file through prior to returning it to the investigator. This is typically a function to reformat data to a more usable format. (**required**)
* **Force root**: A root directory to begin file system search from; this overrides any root specified on the command line. (**optional**)
* **Path filter**: A regular expression to apply to each discovered file, if the pattern does not match the entire path it will not be analyzed further. (**optional**)

Templates can be added to `templates.go`. Currently, there is no method to
supply a more dynamic template from the command line, the templates have
been hardcoded into the module. This is to avoid any potential for sensitive
data being returned from the agent (e.g., password hashes).

