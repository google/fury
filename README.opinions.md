Highly opinionated, probably wrong, if you disagree, you're welcome to
get your own project with its own opinions. Constructive debate
welcome though.

Very short, this is just a skeleton of an essay.

# What's wrong with config management systems

I've used puppet, chef, ansible, fabric, slack (old google thing, not
the chat system). Here's what I think is wrong with the general
direction of config management today.

## Facts are bad

By facts, I mean data snippets gathered from the machine receiving
configuration, e.g. what OS it's running, what its IP address is
(usually derived by some terrible heuristic in everything but
non-trivial cases), and so forth.

Gathering and using these facts encourage poor patterns, complicate
configurations, litter the variable namespace with crap, fail in
surprising ways. Collecting state from a system that is by definition
not in the state you expected it to be should be a deliberate action,
considered carefully by the configuration writer.

## Premature optimization leads to painful DSLs

puppet, chef, ansible all use whacky DSLs to try and abstract away the
fact that imperative programs cannot describe a desired *end state*
cleanly. Some also use the DSL to attempt to optimize the action graph
for better performance.

This optimizes for the wrong thing. Configuration should first be
maintainable and pleasant to write, fast second.

## Modules are needle in a haystack

I need to do about 4 things, but they're not the same 4 things as
everyone else. So, all the config mgmt's modules end up being
ultra-generic to the point where using them requires reading the
manual every time, and the defaults don't do what anyone wants out of
the box.

## Nobody reuses cookbooks/recipes/roles/whatever

All config mgmt systems have a place for people to share their roles,
on the theory that it'll be super easy for new operators, you just
need to apply the 5 roles someone else already wrote and you're done!

If so, why are there 5-50 roles for doing each thing? Because the
roles are either so generic they're irritating to use, or specialized
in just the wrong way. In the end, you end up writing your own anyway.

Until everyone agrees on how to admin systems, roles will not be
reused between unaligned organizations.

## Nobody runs PKIs properly

Stop building systems that rely on agents and a PKI to function. The
three organizations who can operate a PKI correctly are using custom
software anyway, and everyone else just uses them as amazingly
complicated public keys.

Agents are a symptom of the DSL optimization scheme. Once you rid
yourself of the impulse to build the DSLs, the usefulness of an agent
vanishes.

## Just log in as root already

Excluding the case of using automation for specific types of service
deployment, your configuration needs root@machine anyway. Stop
pretending you don't, and don't litter your config mgmt system with
the twelve different mechanisms for privilege escalation that people
think they need. Just bite the bullet and let the config mgmt log in
as root.

Admittedly, this does exclude people who have valid reasons for
restricting privilege escalation. But they're already dealing with
headaches stemming from that anyway, so I'm just going to give them
one more.

## Automation is just software that doesn't have tests

Stop trying to shoehorn sysadmin into non-programming-language config
files. The end result will just be programs written in a bizarre
dialect of something, with no tests.

# What does my ideal config mgmt system look like?

Basically, ansible without the facts, inventory, or yaml. Or like
slack with a language better than shell scripts.

## Agentless

It's easier to make agentless agentful than the other way around, and
agentless works fine for the small/medium deployments I care about.

## Not general purpose

The sins of existing systems all stem from trying to be everything
(case in point, ansible learning to manage switch configuration). By
definition, if you want to do anything, you'll end up with a
turing-complete programming language. Except you built your DSL to be
simpler than a turing-complete programming language, so the end result
will be a really crappy turing-complete programming language.

By giving up on specific aspects of generality, I'm hoping (probably
incorrectly) to get a saner system for the bits I actually do care
about.

## Real programming language

Enough with the DSLs. To a great extent, enough also with library
functions. Accept that people want to build vastly different things,
and give them a proper programming environment to make that
happen. Make it just useful enough that they can build their own
domain-specific helpers quickly, and maintain them easily.

## Cater to the install/configure/run cycle

90% of sysadmin is "install package, write configuration files, start
service". Explicitly cater to that, rather than provide a completely
generic graph workflow system. Provide escape hatches for the special
cases. The escape hatch goes straight to "run arbitrary commands in a
proper programming language", no stopping at intermediate "framework
sub-basement 3 with sprinkles."
