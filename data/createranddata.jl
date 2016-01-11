using Faker

#Create files with random data (ASCII for now) with a specified file name
#The file name is the format: <text><unixtimestamp>.<ext>
#The extension is 'txt' now, but can be changed.

function createfilename()

    filePrefix::AbstractString = "logfile" #Can be changed
    fileTimeStamp = Faker.unix_time()

    fileName = *(filePrefix, "_", string(fileTimeStamp), ".txt")
    return fileName
end

function createRandFile()
    #relPath = "./"
    fileName = createfilename()
    numlines = round(Int, rand() * 1000) #number of lines, 0 to 1000 here

    filehandle = open(fileName, "w+")

    for i in 1:numlines
        write(filehandle, Faker.paragraphs())
    end

    close(filehandle)

end
