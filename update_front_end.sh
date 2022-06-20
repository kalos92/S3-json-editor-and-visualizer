cd frontend/ui
npm run build --build.rollupOptions.external
cd dist
cp -r * ../../../server/build