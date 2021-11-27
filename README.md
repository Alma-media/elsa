# MQTT Routing/Flow manager

### TODO
- Restore state after shutdown

### Config example
```
[
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e0",
        "output": "/trig.in.switch/00"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c360",
        "output": "/trig.in.switch/00"
    },
    {
        "input": "/trig.out.switch/00",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e0"
    },
    {
        "input": "/trig.out.switch/00",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c360"
    },
    {
        "input": "/trig.out.switch/00",
        "output": "/relay/3f8a79dc-854d-49d6-aa24-522422bf9140"
    },
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e1",
        "output": "/trig.in.switch/01"
    },
    {
        "input": "/trig.out.switch/01",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e1"
    },
    {
        "input": "/trig.out.switch/01",
        "output": "/relay/3f8a79dc-854d-49d6-aa24-522422bf9141"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c361",
        "output": "/trig.in.switch/01"
    },
    {
        "input": "/trig.out.switch/01",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c361"
    },
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e2",
        "output": "/trig.in.switch/02"
    },
    {
        "input": "/trig.out.switch/02",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e2"
    },
    {
        "input": "/trig.out.switch/02",
        "output": "/relay/3f8a79dc-854d-49d6-aa24-522422bf9142"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c362",
        "output": "/trig.in.switch/02"
    },
    {
        "input": "/trig.out.switch/02",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c362"
    },
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e3",
        "output": "/trig.in.switch/03"
    },
    {
        "input": "/trig.out.switch/03",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e3"
    },
    {
        "input": "/trig.out.switch/03",
        "output": "/relay/3f8a79dc-854d-49d6-aa24-522422bf9143"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c363",
        "output": "/trig.in.switch/03"
    },
    {
        "input": "/trig.out.switch/03",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c363"
    },
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e4",
        "output": "/trig.in.switch/04"
    },
    {
        "input": "/trig.out.switch/04",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e4"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c364",
        "output": "/trig.in.switch/04"
    },
    {
        "input": "/trig.out.switch/04",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c364"
    },
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e5",
        "output": "/trig.in.switch/05"
    },
    {
        "input": "/trig.out.switch/05",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e5"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c365",
        "output": "/trig.in.switch/05"
    },
    {
        "input": "/trig.out.switch/05",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c365"
    },
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e6",
        "output": "/trig.in.switch/06"
    },
    {
        "input": "/trig.out.switch/06",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e6"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c366",
        "output": "/trig.in.switch/06"
    },
    {
        "input": "/trig.out.switch/06",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c366"
    },
    {
        "input": "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e7",
        "output": "/trig.in.switch/07"
    },
    {
        "input": "/trig.out.switch/07",
        "output": "/led/5d09399e-8a48-41c5-9f47-83127d8a69e7"
    },
    {
        "input": "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c367",
        "output": "/trig.in.switch/07"
    },
    {
        "input": "/trig.out.switch/07",
        "output": "/led/2ccfd6b1-afdb-4c94-a651-5e112084c367"
    }
]
```