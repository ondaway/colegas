package com.ondaway.colegas.person;

public class Person {
    private final String id;
    private final String name;
    private final String surname;
    
    public Person(final String id, final String name, final String surname) {
        this.id      = id;
        this.name    = name;
        this.surname = surname;
    }
    
    public String getId()      { return this.id; }
    public String getName()    { return this.name; }
    public String getSurname() { return this.surname; }
}
