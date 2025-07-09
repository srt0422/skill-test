const dotenv = require('dotenv');
dotenv.config();

const { Pool } = require('pg');
const { processDBRequest } = require('../utils/process-db-request');

// Sample student data
const sampleStudents = [
    {
        name: "Alice Johnson",
        email: "alice.johnson@student.school.com",
        gender: "Female",
        phone: "5551234567",
        dob: "2006-03-15",
        currentAddress: "123 Maple Street, Springfield",
        permanentAddress: "123 Maple Street, Springfield",
        fatherName: "Robert Johnson",
        fatherPhone: "5551234568",
        motherName: "Maria Johnson", 
        motherPhone: "5551234569",
        guardianName: "Robert Johnson",
        guardianPhone: "5551234568",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 10",
        section: "A",
        admissionDate: "2023-09-01",
        roll: 1
    },
    {
        name: "Bob Smith",
        email: "bob.smith@student.school.com",
        gender: "Male",
        phone: "5552345678",
        dob: "2006-07-22",
        currentAddress: "456 Oak Avenue, Springfield",
        permanentAddress: "456 Oak Avenue, Springfield", 
        fatherName: "James Smith",
        fatherPhone: "5552345679",
        motherName: "Linda Smith",
        motherPhone: "5552345680",
        guardianName: "James Smith",
        guardianPhone: "5552345679",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 10",
        section: "A", 
        admissionDate: "2023-09-01",
        roll: 2
    },
    {
        name: "Carol Davis",
        email: "carol.davis@student.school.com",
        gender: "Female",
        phone: "5553456789",
        dob: "2006-11-08",
        currentAddress: "789 Pine Road, Springfield",
        permanentAddress: "789 Pine Road, Springfield",
        fatherName: "Michael Davis",
        fatherPhone: "5553456790",
        motherName: "Sarah Davis",
        motherPhone: "5553456791",
        guardianName: "Michael Davis",
        guardianPhone: "5553456790",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 10",
        section: "B",
        admissionDate: "2023-09-01",
        roll: 3
    },
    {
        name: "David Wilson",
        email: "david.wilson@student.school.com",
        gender: "Male",
        phone: "5554567890",
        dob: "2005-01-14",
        currentAddress: "321 Elm Street, Springfield",
        permanentAddress: "321 Elm Street, Springfield",
        fatherName: "Thomas Wilson",
        fatherPhone: "5554567891",
        motherName: "Jennifer Wilson",
        motherPhone: "5554567892",
        guardianName: "Thomas Wilson",
        guardianPhone: "5554567891",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 11",
        section: "A",
        admissionDate: "2022-09-01",
        roll: 4
    },
    {
        name: "Emma Brown",
        email: "emma.brown@student.school.com",
        gender: "Female",
        phone: "5555678901",
        dob: "2005-05-30",
        currentAddress: "654 Cedar Lane, Springfield",
        permanentAddress: "654 Cedar Lane, Springfield",
        fatherName: "Kevin Brown",
        fatherPhone: "5555678902",
        motherName: "Michelle Brown",
        motherPhone: "5555678903",
        guardianName: "Kevin Brown",
        guardianPhone: "5555678902",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 11",
        section: "A",
        admissionDate: "2022-09-01",
        roll: 5
    },
    {
        name: "Frank Miller",
        email: "frank.miller@student.school.com",
        gender: "Male",
        phone: "5556789012",
        dob: "2005-09-12",
        currentAddress: "987 Birch Avenue, Springfield",
        permanentAddress: "987 Birch Avenue, Springfield",
        fatherName: "Daniel Miller",
        fatherPhone: "5556789013",
        motherName: "Nancy Miller",
        motherPhone: "5556789014",
        guardianName: "Daniel Miller",
        guardianPhone: "5556789013",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 11",
        section: "B",
        admissionDate: "2022-09-01",
        roll: 6
    },
    {
        name: "Grace Lee",
        email: "grace.lee@student.school.com",
        gender: "Female",
        phone: "5557890123",
        dob: "2004-12-25",
        currentAddress: "147 Spruce Drive, Springfield",
        permanentAddress: "147 Spruce Drive, Springfield",
        fatherName: "Steven Lee",
        fatherPhone: "5557890124",
        motherName: "Karen Lee",
        motherPhone: "5557890125",
        guardianName: "Steven Lee",
        guardianPhone: "5557890124",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 12",
        section: "A",
        admissionDate: "2021-09-01",
        roll: 7
    },
    {
        name: "Henry Taylor",
        email: "henry.taylor@student.school.com",
        gender: "Male",
        phone: "5558901234",
        dob: "2004-04-18",
        currentAddress: "258 Walnut Street, Springfield",
        permanentAddress: "258 Walnut Street, Springfield",
        fatherName: "Paul Taylor",
        fatherPhone: "5558901235",
        motherName: "Rebecca Taylor",
        motherPhone: "5558901236",
        guardianName: "Paul Taylor",
        guardianPhone: "5558901235",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 12",
        section: "A",
        admissionDate: "2021-09-01",
        roll: 8
    },
    {
        name: "Iris Anderson",
        email: "iris.anderson@student.school.com",
        gender: "Female",
        phone: "5559012345",
        dob: "2004-08-03",
        currentAddress: "369 Chestnut Road, Springfield",
        permanentAddress: "369 Chestnut Road, Springfield",
        fatherName: "Mark Anderson",
        fatherPhone: "5559012346",
        motherName: "Lisa Anderson",
        motherPhone: "5559012347",
        guardianName: "Mark Anderson",
        guardianPhone: "5559012346",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 12",
        section: "B",
        admissionDate: "2021-09-01",
        roll: 9
    },
    {
        name: "Jack Thompson",
        email: "jack.thompson@student.school.com",
        gender: "Male",
        phone: "5550123456",
        dob: "2004-10-27",
        currentAddress: "741 Poplar Avenue, Springfield",
        permanentAddress: "741 Poplar Avenue, Springfield",
        fatherName: "Gary Thompson",
        fatherPhone: "5550123457",
        motherName: "Helen Thompson",
        motherPhone: "5550123458",
        guardianName: "Gary Thompson",
        guardianPhone: "5550123457",
        relationOfGuardian: "Father",
        systemAccess: true,
        class: "Grade 12",
        section: "B",
        admissionDate: "2021-09-01",
        roll: 10
    }
];

async function seedDatabase() {
    console.log('🌱 Starting database seeding...');
    
    try {
        // Create classes first
        console.log('📚 Creating classes...');
        const classes = ['Grade 10', 'Grade 11', 'Grade 12'];
        for (const className of classes) {
            const insertClassQuery = 'INSERT INTO classes (name) VALUES ($1) ON CONFLICT (name) DO NOTHING';
            await processDBRequest({ query: insertClassQuery, queryParams: [className] });
            console.log(`   ✅ Created class: ${className}`);
        }

        // Create sections
        console.log('📝 Creating sections...');
        const sections = ['A', 'B'];
        for (const sectionName of sections) {
            const insertSectionQuery = 'INSERT INTO sections (name) VALUES ($1) ON CONFLICT (name) DO NOTHING';
            await processDBRequest({ query: insertSectionQuery, queryParams: [sectionName] });
            console.log(`   ✅ Created section: ${sectionName}`);
        }

        // Create students using the stored procedure
        console.log('👥 Creating students...');
        for (let i = 0; i < sampleStudents.length; i++) {
            const student = sampleStudents[i];
            const query = 'SELECT * FROM student_add_update($1)';
            const queryParams = [JSON.stringify(student)];
            
            const result = await processDBRequest({ query, queryParams });
            const { userId, status, message } = result.rows[0];
            
            if (status) {
                console.log(`   ✅ Created student #${i + 1}: ${student.name} (ID: ${userId})`);
            } else {
                console.log(`   ❌ Failed to create student #${i + 1}: ${student.name} - ${message}`);
            }
        }

        console.log('🎉 Database seeding completed successfully!');
        console.log(`📊 Created ${classes.length} classes, ${sections.length} sections, and ${sampleStudents.length} students`);
        
    } catch (error) {
        console.error('❌ Error seeding database:', error);
        throw error;
    }
}

// Export the function for use in other scripts
module.exports = { seedDatabase };

// Run the seeding if this file is executed directly
if (require.main === module) {
    seedDatabase()
        .then(() => {
            console.log('✨ Seeding process completed!');
            process.exit(0);
        })
        .catch((error) => {
            console.error('💥 Seeding process failed:', error);
            process.exit(1);
        });
}
