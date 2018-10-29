# Stepmania Simfile Format
Stepmania uses a custom format for representing Chart data. This format contains metadata about the Song and the associated Chart(s).

## Format Support
This parser currently supports only the `.sm` extension. See the [Stepmania Wiki](https://github.com/stepmania/stepmania/wiki/sm) for details.

## Simfile Structure
A Simfile consists of two major sections: 
* **Header Tags**
* **Chart Data**

### Header Tags
The **Header** section contains metadata about a Simfile. 
> NOTE: Some tags contain lists of beat/value pairs. These must be parsed.

### Chart Data
A **Chart** contains metadata for the specified collection of notes. 
> NOTE: Most charts do not contain values for the Groove Radar.