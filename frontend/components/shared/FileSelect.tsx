import { UploadIcon } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import AddFileForm from "@/app/data/AddFileForm";
import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { listFiles } from "@/lib/redux/slices/fileListSlice/asyncActions";
import { selectFileList, selectFileListPagination } from "@/lib/redux/slices/fileListSlice/selectors";
import { fileListThunk } from "@/lib/redux/slices/fileListSlice/thunks";
import { cn } from "@/lib/utils";

interface File {
  ID: string;
  Filename: string;
  S3URI: string;
}

interface FileFilters {
  filename?: string;
  id?: string;
}

interface FileSelectProps {
  onChange: (value: string) => void;
  value: string;
  label: string;
  globPatterns?: string[];
}

export function FileSelect({ onChange, value, label, globPatterns }: FileSelectProps) {
  const dispatch = useDispatch();
  const files = useSelector(selectFileList);
  const pagination = useSelector(selectFileListPagination);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [open, setOpen] = useState<boolean>(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchTerm, setSearchTerm] = useState("");

  const cidPattern = /^[a-zA-Z0-9]{46}$/;

  let filenameFilter = globPatterns?.map((pattern) => pattern.replace("*", "%")).join("|") || "%";
  let fileAcceptExtensions = globPatterns?.map((pattern) => pattern.replace("*", "")).join(",") || "*";
  let filters: FileFilters = useMemo(() => ({}), []);

  if (searchTerm && cidPattern.test(searchTerm)) {
    filters.id = searchTerm;
    filters.filename = filenameFilter;
  } else {
    filters.filename = searchTerm ? `%${searchTerm}%` : filenameFilter;
  }

  useEffect(() => {
    setLoading(true);
    setError(null);
    // @ts-ignore
    dispatch(fileListThunk({ page: currentPage, pageSize: 10, filters }))
      .unwrap()
      .then(() => setLoading(false))
      .catch((fetchError: Error) => {
        setError(fetchError.message);
        setLoading(false);
      });
  }, [dispatch, currentPage, filters, searchTerm, open]);

  const getFileValue = (file: File): string => `${file?.S3URI}`;

  const handleUpload = async (id: string) => {
    const files = await listFiles({ page: 1, pageSize: 10, filters: { id: id, filename: filenameFilter } });
    if (files?.data?.length && files?.data[0]?.ID === id) {
      onChange(getFileValue(files.data[0]));  // This will update the form value
    }
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <div className="relative flex items-center gap-2">
        <PopoverTrigger asChild>
          <Button className="justify-between p-3 text-base grow" variant="input" role="combobox" aria-expanded={open}>
            <span className={cn(!value && "text-muted-foreground")}>{value ? value.split("/").pop() : `Select ${label} file...`}</span>
          </Button>
        </PopoverTrigger>
        <AddFileForm
          trigger={
            <Button variant="secondary" size="xs" className="absolute z-30 right-5 group hover:bg-secondary">
              <span className="w-0 overflow-hidden group-hover:w-12 transition-[width]">Upload</span>
              <UploadIcon className="w-4 h-4 ml-2" />
            </Button>
          }
          accept={fileAcceptExtensions}
          onUpload={handleUpload}
        />
      </div>
      <PopoverContent className="p-0 " style={{ width: `var(--radix-popover-trigger-width)` }}>
        <Command>
          <CommandInput placeholder="Search files..." onValueChange={setSearchTerm} value={searchTerm} />
          <CommandEmpty>No file found.</CommandEmpty>
          {error && <CommandItem>Error fetching files.</CommandItem>}
          <CommandGroup>
            {files.map((file) => (
              <CommandItem
                key={file.ID}
                value={getFileValue(file)}
                onSelect={() => {
                  onChange(getFileValue(file) === value ? "" : getFileValue(file));
                  setOpen(false);
                }}
              >
                {file.Filename}
              </CommandItem>
            ))}
          </CommandGroup>
        </Command>
      </PopoverContent>
    </Popover>
  );
}