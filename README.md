Fail2Connect - Ban connections that fail to connect
===========================================================
Fail2Connect is a program written in Golang. It reads log files like /var/log/openvpn.log or /var/log/mail.log. Each time
a new line is added, Fail2Connect will try to match the new line against two regular expressions. These regular expressions are
defined in the config.json file. The first regular expression will match against a line that mentions a new connection from a certain IP address is made.
This doesn't mean that that connection is from a valid user. It only means that someone made a connection. The second regular expression
will match against a line that indicates that a connection with a certain IP address is successfully made, so from a valid user.

Example from OpenVPN:

``TCP connection established with [AF_INET]46.143.188.17:43912 (will match the first regular expression)
...
...
...
Peer Connection Initiated with [AF_INET]46.143.188.17:43912 (will match the second regular expression)``

The first line mentions a new connection is made from some IP address, the second line mentions that that IP address
successfully made a connection.

Fail2Connect will track the time between the first line and the second line. If it takes longer than defined in the
config file, it will run a command defined in the config file. This command should ban that IP address. This can be done
because only connections from valid users will match the second line. If a connection from an invalid user is made, it will not
match the second line.


## Configuration
--------------------------------------------------------
The config.json file is a JSON file containing multiple watchers. Each watcher is a JSON object containing the following
keys:

- enabled (should this watcher be enabled? true or false)
- path_to_log_file (location of the log file)
- connection_regex (egular expression for matching incoming, but not yet valid, connections. MUST HAVE A CAPTURE GROUP FOR THE IP ADDRESS)
- success_regex (regular expression for matching connections from valid users. MUST HAVE A CAPTURE GROUP FOR THE IP ADDRESS)
- ban_command (command to execute when there is too much time between the connection line and success line. IP_TO_BAN will be replaced with the corresponding IP address)
- ultimatum_time_in_seconds (maximum time between connection and success line)
- trust_known (should known IP addresses be trusted? true or false. This means that every IP address from the success regular expression will be trusted, meaning those IPs won't be banned even if it does not reach the success line next time)
- instant_ban_after (in case of a DOS attack, after how many attempts within the ultimatum should the IP address be banned?)