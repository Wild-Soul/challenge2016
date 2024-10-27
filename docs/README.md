# Distributor Permissions System

Thank you for the opportunity to complete this takehome assignment. I've implemented a geographical distribution rights management system in Go that handles permissions for distributors in a hierarchical structure.

## Problem Understanding

The system addresses the following requirements:
1. Managing hierarchical distributor relationships (parent-child)
2. Handling geographical permissions at country, state/province, and city levels
3. Processing both inclusion and exclusion rules
4. Validating permissions against parent distributors
5. Efficiently loading and processing geographical data

## Solution Design

### Permission Management

I implemented a hierarchical approach to manage permissions:

```
Parent Distributor
    └── Child Distributor
         └── Grandchild Distributor
```

The system follows these key rules:
- Child distributors can only have permissions that are a subset of their parent's permissions
- Permissions are explicitly granted through INCLUDE statements
- Restrictions are applied through EXCLUDE statements
- The same location cannot be both included and excluded
- By default, distributors have no permissions until explicitly granted

### Core Algorithm

The permission checking algorithm follows this sequence:

1. **Parent Validation**
   - First checks parent distributor's permissions if present
   - Automatically denies if parent doesn't have permission

2. **Exclusion Processing**
   - Checks for matching excluded patterns
   - Supports patterns at country, state, or city level
   - Denies permission on match

3. **Inclusion Processing**
   - Checks for matching included patterns
   - Supports patterns at country, state, or city level
   - Grants permission on match

4. **Default Behavior**
   - Denies permission if no matches found (secure by default)

### Location Pattern Matching

I implemented a hierarchical pattern matching system:
- Full location format: CITY-PROVINCE-COUNTRY
- Supports partial matching at any level: COUNTRY, PROVINCE-COUNTRY, or CITY-PROVINCE-COUNTRY
- Uses right-to-left matching for partial patterns

## Implementation Details

### Key Data Structures

1. **Location**
```go
type Location struct {
    City     string
    Province string
    Country  string
}
```

2. **Distributor**
```go
type Distributor struct {
    Name     string
    Parent   *Distributor
    Includes map[string]bool
    Excludes map[string]bool
    mu       sync.RWMutex
}
```

### Performance Considerations

1. **Concurrent Data Processing**
   - Implemented concurrent CSV loading using goroutines
   - Used WaitGroup for synchronization
   - Ensured thread-safe operations

2. **Optimized Lookups**
   - Utilized maps for O(1) permission checks
   - Implemented hierarchical location data caching
   - Added mutex protection for thread safety

## Running the Application

### Command Line Interface

```bash
make run
```

Command line flags used in run stage of makefile:
- `-csv`: Location database file path
- `-perm`: Permissions configuration file path
- `-dist`: Target distributor name
- `-check`: Location to verify permission for

### Input File Formats

1. **Locations Database (CSV)**
```csv
City Code,Province Code,Country Code,City Name,Province Name,Country Name
PUNCH,JK,IN,Punch,Jammu and Kashmir,India
```

2. **Permissions Configuration**
```
Permissions for DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN
```

## Testing

- Yet to add tests.

## Trade-offs and Future Improvements

1. **Memory Usage**
   - Current implementation keeps location database in memory
   - Could be optimized with database backend for larger datasets, but task strictly specified not to do so.

2. **Concurrency Handling**
   - Implemented basic thread safety
   - Could be enhanced with connection pooling for higher throughput

3. **Storage Layer**
   - Direct memory implementation for simplicity
   - Abstracted storage interface allows for easy database integration

## Questions?

Please don't hesitate to reach out if you have any questions about my implementation or would like to discuss any aspects of the solution in more detail.
