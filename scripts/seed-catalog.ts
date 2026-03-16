import mongoose from 'mongoose';
import path from 'path';
import { readFileSync } from 'fs';

const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://root:12345abc@10.2.20.113:27017/calculadora_paneles?authSource=admin&directConnection=true';

async function seed() {
  await mongoose.connect(MONGODB_URI);
  console.log('Connected to MongoDB');

  const db = mongoose.connection.db!;

  // Seed panels
  const panelsData = JSON.parse(
    readFileSync(path.join(__dirname, '../backend/src/data/default-panels.json'), 'utf-8')
  );

  const panelsCollection = db.collection('panelcatalogs');
  await panelsCollection.deleteMany({});
  const panelsWithDefaults = panelsData.map((p: any) => ({ ...p, isActive: true }));
  await panelsCollection.insertMany(panelsWithDefaults);
  console.log(`✓ Seeded ${panelsData.length} panels`);

  // Seed inverters
  const invertersData = JSON.parse(
    readFileSync(path.join(__dirname, '../backend/src/data/default-inverters.json'), 'utf-8')
  );

  const invertersCollection = db.collection('invertercatalogs');
  await invertersCollection.deleteMany({});
  const invertersWithDefaults = invertersData.map((i: any) => ({ ...i, isActive: true }));
  await invertersCollection.insertMany(invertersWithDefaults);
  console.log(`✓ Seeded ${invertersData.length} inverters`);

  await mongoose.disconnect();
  console.log('✓ Seed complete');
}

seed().catch((err) => {
  console.error('Seed failed:', err);
  process.exit(1);
});
